package task

type Slicer []*Task

func (s Slicer) TaskGate(gat *Gate) Slicer {
	var lis Slicer

	for _, x := range s {
		if x.Has(&Task{Gate: gat}) {
			lis = append(lis, x)
		}
	}

	return lis
}

func (s Slicer) TaskMeta(met *Meta) Slicer {
	var lis Slicer

	for _, x := range s {
		if x.Has(&Task{Meta: met}) {
			lis = append(lis, x)
		}
	}

	return lis
}
