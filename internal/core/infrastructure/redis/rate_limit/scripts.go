package core_redis_rate_limit

const luaScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])

local current = redis.call('GET', key)

if current == false then
	redis.call('SET', key, 1)
	redis.call('EXPIRE', key, window)

	return {1, limit - 1, window}
end

local count = tonumber(current)

if count >= limit then
	local ttl = redis.call('TTL', key)

	return {0, 0, ttl}
end

local new_count = redis.call('INCR', key)
local ttl = redis.call('TTL', key)

if ttl == -1 then
	redis.call('EXPIRE', key, window)
	ttl = window
end

return {1, limit - new_count, ttl}
`
