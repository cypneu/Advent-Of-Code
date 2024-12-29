package main

import "fmt"

func Solution9A(input string) {
	diskmap := make([]int, len(input))
	for i, r := range input {
		diskmap[i] = int(r - '0')
	}

	total, block := 0, 0
	right := len(diskmap) - 1
	if right%2 == 1 {
		right--
	}

	left := 0
	for left <= right {
		if left%2 == 0 {
			fileId := left / 2
			total += fileId * (2*block + diskmap[left] - 1) * (diskmap[left]) / 2
			block += diskmap[left]
			left++
		} else if diskmap[left] <= diskmap[right] {
			diskmap[right] -= diskmap[left]
			fileId := right / 2
			total += fileId * (2*block + diskmap[left] - 1) * (diskmap[left]) / 2
			block += diskmap[left]
			left++
		} else {
			diskmap[left] -= diskmap[right]
			fileId := right / 2
			total += fileId * (2*block + diskmap[right] - 1) * (diskmap[right]) / 2
			block += diskmap[right]
			right -= 2
		}
	}
	fmt.Println(total)
}

func Solution9B(input string) {
	diskmap := make([]int, len(input))
	indicesTotal := make([]int, len(diskmap))
	for i, r := range input {
		diskmap[i] = int(r - '0')
		if i > 0 {
			indicesTotal[i] += diskmap[i-1] + indicesTotal[i-1]
		}
	}

	total := 0
	right := len(diskmap) - 1
	if right%2 == 1 {
		right--
	}

	for right > 0 {
		left := 1
		for left < right {
			if diskmap[left] >= diskmap[right] {
				fileId := right / 2
				total += fileId * (2*indicesTotal[left] + diskmap[right] - 1) * (diskmap[right]) / 2
				diskmap[left] -= diskmap[right]
				indicesTotal[left] += diskmap[right]
				diskmap[right] = 0
			}

			left += 2
		}

		right -= 2
	}

	for i := 0; i < len(diskmap); i += 2 {
		fileId := i / 2
		total += fileId * (2*indicesTotal[i] + diskmap[i] - 1) * (diskmap[i]) / 2
	}

	fmt.Println(total)
}
