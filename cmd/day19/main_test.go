package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	example := []struct {
		name  string
		in    string
		boxes []block
	}{
		{
			"first example",
			`--- scanner 0 ---
0,2,0
4,1,0
3,3,0

--- scanner 1 ---
-1,-1,0
-5,0,0
-2,1,0
`,
			[]block{
				{
					Sensors: map[int]v.I3{0: {}},
					Beacons: []v.I3{
						{X: 0, Y: 2},
						{X: 3, Y: 3},
						{X: 4, Y: 1},
					},
				},
				{
					Sensors: map[int]v.I3{1: {}},
					Beacons: []v.I3{
						{X: -5, Y: 0},
						{X: -2, Y: 1},
						{X: -1, Y: -1},
					},
				},
			},
		},
	}
	for _, tc := range example {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := read(strings.NewReader(tc.in))
			require.NoError(t, err)
			for i, box := range got {
				sort.Sort(byXYZ(box.Beacons))
				assert.Equal(t, tc.boxes[i], box)
			}
		})
	}
}

func TestBeaconsMatch(t *testing.T) {
	var (
		box1 = block{
			Sensors: map[int]v.I3{0: {}},
			Beacons: []v.I3{
				{X: 0, Y: 2},
				{X: 3, Y: 3},
				{X: 4, Y: 1},
			},
		}

		box2 = block{
			Sensors: map[int]v.I3{1: {}},
			Beacons: []v.I3{
				{X: -5, Y: 0},
				{X: -2, Y: 1},
				{X: -1, Y: -1},
			},
		}

		extraBeacon = block{
			Sensors: map[int]v.I3{2: {}},
			Beacons: append(box1.Beacons, v.I3{X: -1, Y: -1, Z: -1}),
		}
	)
	sort.Sort(byXYZ(extraBeacon.Beacons))

	tt := []struct {
		name       string
		a, b       block
		minQty     int
		wantOffset v.I3
		wantMatch  bool
	}{

		{
			name:       "example 0",
			a:          box1,
			b:          box2,
			minQty:     3,
			wantOffset: v.I3{X: 5, Y: 2},
			wantMatch:  true,
		},
		{
			name:      "insufficient matches => false",
			a:         box1,
			b:         box2,
			minQty:    4,
			wantMatch: false,
		},
		{
			name:       "an extra beacon in a => true",
			a:          extraBeacon,
			b:          box2,
			minQty:     3,
			wantOffset: v.I3{X: 5, Y: 2},
			wantMatch:  true,
		},
		{
			name:       "an extra beacon in b => true",
			a:          box2,
			b:          extraBeacon,
			minQty:     3,
			wantOffset: v.I3{X: -5, Y: -2},
			wantMatch:  true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			offset, isMatch := beaconsMatch(tc.a.Beacons, tc.b.Beacons, tc.minQty)
			assert.Equal(t, tc.wantOffset, offset)
			assert.Equal(t, tc.wantMatch, isMatch)
		})
	}
}

func TestRotations(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	unique := make(map[v.I3]int, 24)
	vect := v.I3{X: 2, Y: 3, Z: 5}

	for i, rot := range getRotations() {
		det, err := rot.Determinant()
		r.NoError(err)
		r.Equal(1, det, "i=%d: rotations must preserve scale and signedness", i)

		vect2, err := vect.Transform(rot)
		r.NoError(err)

		if prev, exists := unique[vect2]; exists {
			t.Logf("\nalready found rotation %d:\n\t%v\n\t%v\n\t%v",
				i, rot[0], rot[1], rot[2])
			t.Logf("previously found at i: %d", prev)
			t.Fail()
		}

		unique[vect2] = i
	}

	a.Equal(24, len(unique))
}

func TestRotateToMatch(t *testing.T) {
	tt := []struct {
		name   string
		minQty int
		in     string
	}{
		{
			"same sensor location", 5,
			`--- scanner 0 ---
0,0,0
-1,-1,1
-2,-2,2
-3,-3,3
-2,-3,1
5,6,-4
8,0,7

--- scanner 0 ---
1,-1,1
4,4,4
2,-2,2
3,-3,3
2,-1,3
-5,4,-6
-8,-7,0

--- scanner 0 ---
-1,-1,-1
-2,-2,-2
3,11,5
-3,-3,-3
-1,-3,-2
4,6,5
-7,0,8

--- scanner 0 ---
1,1,-1
9,8,7
2,2,-2
3,3,-3
1,3,-2
-4,-6,5
7,0,8

--- scanner 0 ---
1,1,1
2,2,2
3,3,3
3,1,2
-6,-4,-5
0,7,-8

--- scanner 0 ---
2,2,2
3,3,3
4,4,4
4,2,3
-5,-4,-3
1,8,-7
`,
		},
		{
			"different sensor locations", 12,
			`--- scanner 0 ---
404,-588,-901
528,-643,409
-838,591,734
390,-675,-793
-537,-823,-458
-485,-357,347
-345,-311,381
-661,-816,-575
-876,649,763
-618,-824,-621
553,345,-567
474,580,667
-447,-329,318
-584,868,-557
544,-627,-890
564,392,-477
455,729,728
-892,524,684
-689,845,-530
423,-701,434
7,-33,-71
630,319,-379
443,580,662
-789,900,-551
459,-707,401

--- scanner 1 ---
686,422,578
605,423,415
515,917,-361
-336,658,858
95,138,22
-476,619,847
-340,-569,-846
567,-361,727
-460,603,-452
669,-402,600
729,430,532
-500,-761,534
-322,571,750
-466,-666,-811
-429,-592,574
-355,545,-477
703,-491,-529
-328,-685,520
413,935,-424
-391,539,-444
586,-435,557
-364,-763,-893
807,-499,-711
755,-354,-619
553,889,-390
`,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			boxes, err := read(strings.NewReader(tc.in))
			require.NoError(t, err)

			for i := 0; i < len(boxes)-1; i++ {
				for j := i + 1; j < len(boxes); j++ {
					t.Run(fmt.Sprintf("%d_%d", i, j), func(t *testing.T) {
						a := assert.New(t)
						box2, ok := rotateToMatch(boxes[i], boxes[j], tc.minQty)
						a.Truef(ok, "boxes %d and %d should match", i, j)
						if ok {
							t.Log("box1", boxes[i])
							t.Log("box2", box2)
						}
					})
				}
			}
		})
	}
}

func TestPart1(t *testing.T) {
	boxes, err := read(strings.NewReader(`--- scanner 0 ---
404,-588,-901
528,-643,409
-838,591,734
390,-675,-793
-537,-823,-458
-485,-357,347
-345,-311,381
-661,-816,-575
-876,649,763
-618,-824,-621
553,345,-567
474,580,667
-447,-329,318
-584,868,-557
544,-627,-890
564,392,-477
455,729,728
-892,524,684
-689,845,-530
423,-701,434
7,-33,-71
630,319,-379
443,580,662
-789,900,-551
459,-707,401

--- scanner 1 ---
686,422,578
605,423,415
515,917,-361
-336,658,858
95,138,22
-476,619,847
-340,-569,-846
567,-361,727
-460,603,-452
669,-402,600
729,430,532
-500,-761,534
-322,571,750
-466,-666,-811
-429,-592,574
-355,545,-477
703,-491,-529
-328,-685,520
413,935,-424
-391,539,-444
586,-435,557
-364,-763,-893
807,-499,-711
755,-354,-619
553,889,-390

--- scanner 2 ---
649,640,665
682,-795,504
-784,533,-524
-644,584,-595
-588,-843,648
-30,6,44
-674,560,763
500,723,-460
609,671,-379
-555,-800,653
-675,-892,-343
697,-426,-610
578,704,681
493,664,-388
-671,-858,530
-667,343,800
571,-461,-707
-138,-166,112
-889,563,-600
646,-828,498
640,759,510
-630,509,768
-681,-892,-333
673,-379,-804
-742,-814,-386
577,-820,562

--- scanner 3 ---
-589,542,597
605,-692,669
-500,565,-823
-660,373,557
-458,-679,-417
-488,449,543
-626,468,-788
338,-750,-386
528,-832,-391
562,-778,733
-938,-730,414
543,643,-506
-524,371,-870
407,773,750
-104,29,83
378,-903,-323
-778,-728,485
426,699,580
-438,-605,-362
-469,-447,-387
509,732,623
647,635,-688
-868,-804,481
614,-800,639
595,780,-596

--- scanner 4 ---
727,592,562
-293,-554,779
441,611,-461
-714,465,-776
-743,427,-804
-660,-479,-426
832,-632,460
927,-485,-438
408,393,-506
466,436,-512
110,16,151
-258,-428,682
-393,719,612
-211,-452,876
808,-476,-593
-575,615,604
-485,667,467
-680,325,-822
-627,-443,-432
872,-547,-609
833,512,582
807,604,487
839,-516,451
891,-625,532
-652,-548,-490
30,-46,-14
`))

	require.NoError(t, err)

	got := part1(boxes)
	assert.Equal(t, 79, got)
}
