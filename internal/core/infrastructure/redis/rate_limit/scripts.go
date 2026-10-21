package core_redis_rate_limit

const luaScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window_ms = tonumber(ARGV[2])

local current = redis.call('GET', key)

if not current then
    redis.call('SET', key, 1, 'PX', window_ms)
    return {1, limit - 1, window_ms}
end

local count = tonumber(current)
if not count then
    redis.call('SET', key, 1, 'PX', window_ms)
    return {1, limit - 1, window_ms}
end

local ttl = redis.call('PTTL', key)
if ttl < 0 then
    redis.call('PEXPIRE', key, window_ms)
    ttl = window_ms
end

if count >= limit then
    return {0, 0, ttl}
end

local new_count = redis.call('INCR', key)
ttl = redis.call('PTTL', key)
if ttl < 0 then
    redis.call('PEXPIRE', key, window_ms)
    ttl = window_ms
end

local remaining = limit - new_count
if remaining < 0 then
    remaining = 0
end

return {1, remaining, ttl}
`
