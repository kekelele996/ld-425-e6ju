import { useCallback, useState } from 'react'

// 通用分页 hook。
export function usePagination(initialPage = 1, initialPageSize = 10) {
  const [page, setPage] = useState(initialPage)
  const [pageSize, setPageSize] = useState(initialPageSize)

  const onChange = useCallback((nextPage: number, nextPageSize: number) => {
    setPage(nextPage)
    setPageSize(nextPageSize)
  }, [])

  return { page, pageSize, onChange }
}
