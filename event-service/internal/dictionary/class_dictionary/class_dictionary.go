package class_dictionary

import "main/internal/models/event"

var ClassDictionaryCategory = map[int][]event.ClassCategoryType{
	9:  {event.Class9, event.Class9_10, event.Class9_11},
	10: {event.Class10, event.Class10_11, event.Class9_10, event.Class9_11},
	11: {event.Class10_11, event.Class11, event.Class9_11},
}

var AllClassCategory = []event.ClassCategoryType{
	event.Class9, event.Class9_10, event.Class9_11,
	event.Class10, event.Class10_11, event.Class11,
}
