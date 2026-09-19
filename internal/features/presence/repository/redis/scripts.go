package presence_redis_repository

const registerScript = `
local key = KEYS[1]

local now = ARGV[1]
local expires = ARGV[2]
local member = ARGV[3]

-- Удаляем протухшие websocket connections.
redis.call(
	"ZREMRANGEBYSCORE",
	key,
	"-inf",
	now
)

-- Был ли пользователь online до регистрации этого connection?
local was_online =
	redis.call("ZCARD", key) > 0

redis.call(
	"ZADD",
	key,
	expires,
	member
)

if was_online then
	return 0
end

return 1
`

const heartbeatScript = `
local key = KEYS[1]

local now = ARGV[1]
local expires = ARGV[2]
local member = ARGV[3]

redis.call(
	"ZREMRANGEBYSCORE",
	key,
	"-inf",
	now
)

if redis.call("ZSCORE", key, member) == false then
	return 0
end

redis.call(
	"ZADD",
	key,
	expires,
	member
)

return 1
`

const unregisterScript = `
local key = KEYS[1]

local now = ARGV[1]
local session_id = ARGV[2]
local connection_id = ARGV[3]

local member = session_id .. ":" .. connection_id

redis.call(
	"ZREM",
	key,
	member
)

redis.call(
	"ZREMRANGEBYSCORE",
	key,
	"-inf",
	now
)

local remaining_session = 0

local members = redis.call(
	"ZRANGE",
	key,
	0,
	-1
)

local prefix = session_id .. ":"

for _, member in ipairs(members) do
	if string.sub(member, 1, string.len(prefix)) == prefix then
		remaining_session = remaining_session + 1
	end
end

local remaining_user = redis.call(
	"ZCARD",
	key
)

return {
	remaining_session,
	remaining_user
}
`

const getOnlineSessionsScript = `
local key = KEYS[1]

local now = ARGV[1]

redis.call(
	"ZREMRANGEBYSCORE",
	key,
	"-inf",
	now
)

return redis.call(
	"ZRANGE",
	key,
	0,
	-1
)
`
