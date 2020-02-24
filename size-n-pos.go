package main

import (
	"math"
)

type box_type struct {
	W int `json:"w"`
	H int `json:"h"`
	X int `json:"x"`
	Y int `json:"y"`
}

type tweet_box_type struct {
	W int `json:"w"`
	H int `json:"h"`
	X int `json:"x"`
	Y int `json:"y"`
	Name string `json:"name"`
	URL string `json:"url"`
	TweetVolume int64 `json:"tweet_volume"`
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func trendsToBoxes(trends []custom_trend, max_width int) []tweet_box_type {
	var ratio float64

	var maxTwVol int
	for _, t := range trends {
		maxTwVol = max(int(t.TweetVolume), maxTwVol)
	}

	ratio = float64(max_width) / float64(maxTwVol)
	t_boxes := make([]tweet_box_type, 0)

	for _, t := range trends {
		width := math.Max(float64(t.TweetVolume) * ratio, 200)
		height := math.Max(width/4.0, 200)
		t_boxes = append(
			t_boxes,
			tweet_box_type{
				int(width), int(height), 0, 0, t.Name, t.URL, t.TweetVolume,
			},
		)
	}

	return t_boxes
}

// Credits to https://github.com/mapbox/potpack
func potpack(boxes []tweet_box_type) {
    // calculate total boxes[j] area and maximum boxes[j] width
    var area float64
    var maxWidth float64

    for _, box := range boxes {
        area += float64(box.W * box.H)
        maxWidth = math.Max(maxWidth, float64(box.W))
    }

    // aim for a squarish resulting container,
    // slightly adjusted for sub-100% space utilization
    var startWidth = math.Max(math.Ceil(math.Sqrt(area / 0.95)), maxWidth)

    // start with a single empty space, unbounded at the bottom
	var spaces = make([]box_type, 0)
	spaces = append(spaces, box_type{int(startWidth), 2147483647, 0, 0})

    var width = 0
    var height = 0

    for j := 0; j <= len(boxes) - 1; j++ {
        // look through spaces backwards so that we check smaller spaces first
        for i := len(spaces) - 1; i >= 0; i-- {
            // look for empty spaces that can accommodate the current box
            if (boxes[j].W > spaces[i].W || boxes[j].H > spaces[i].H) { continue }

            // found the space; add the box to its top-left corner
            // |-------|-------|
            // |  box  |       |
            // |_______|       |
            // |         space |
            // |_______________|
            boxes[j].X = spaces[i].X
			boxes[j].Y = spaces[i].Y

            height = max(height, boxes[j].Y + boxes[j].H)
            width = max(width, boxes[j].X + boxes[j].W)

            if boxes[j].W == spaces[i].W && boxes[j].H == spaces[i].H {
				// space matches the box exactly; remove it
				var last box_type
				last, spaces = spaces[len(spaces)-1], spaces[:len(spaces)-1]

                if i < len(spaces) { spaces[i] = last }

            } else if boxes[j].H == spaces[i].H {
                // space matches the box height; update it accordingly
                // |-------|---------------|
                // |  box  | updated space |
                // |_______|_______________|
                spaces[i].X += boxes[j].W;
                spaces[i].W -= boxes[j].W;

            } else if boxes[j].W == spaces[i].W {
                // space matches the box width; update it accordingly
                // |---------------|
                // |      box      |
                // |_______________|
                // | updated space |
                // |_______________|
                spaces[i].Y += boxes[j].H;
                spaces[i].H -= boxes[j].H;

            } else {
                // otherwise the box splits the space into two spaces
                // |-------|-----------|
                // |  box  | new space |
                // |_______|___________|
                // | updated space     |
				// |___________________|
				spaces = append(
					spaces,
					box_type{spaces[i].W - boxes[j].W, boxes[j].H, spaces[i].X + boxes[j].W, spaces[i].Y},
				)

				spaces[i].Y += boxes[j].H
				spaces[i].H -= boxes[j].H
			}
            break;
        }
    }
}