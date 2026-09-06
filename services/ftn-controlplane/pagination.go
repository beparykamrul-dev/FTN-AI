package controlplane
func PageLimit(limit,defaultLimit,maxLimit int)int{if maxLimit<1{return 0};if defaultLimit<1{defaultLimit=1};if defaultLimit>maxLimit{defaultLimit=maxLimit};if limit<1{return defaultLimit};if limit>maxLimit{return maxLimit};return limit}
