package gateway

import "strconv"

// parseID 解析 URL 里的数字 ID
func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
