package pkg

func GetPaginationIndexes(pageNumber, itemsPerPage int) (startIndex, endIndex int) {
	startIndex = (pageNumber - 1) * itemsPerPage
	endIndex = int(startIndex + itemsPerPage - 1)
	return
}
