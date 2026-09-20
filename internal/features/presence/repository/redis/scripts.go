package presence_redis_repository

const registerScript = `
local key = KEYS[1]

local now = ARGV[1]
local expires = ARGV[2]
local member = ARGV[3]
local session_prefix = ARGV[4]

redis.call(
	"ZREMRANGEBYSCORE",
	key,
	"-inf",
	now
)

local was_user_online =
	redis.call("ZCARD", key) > 0

local was_session_online = 0

local members = redis.call(
	"ZRANGE",
	key,
	0,
	-1
)

for _, current_member in ipairs(members) do
	if string.sub(
		current_member,
		1,
		string.len(session_prefix)
	) == session_prefix then
		was_session_online = 1
		break
	end
end

redis.call(
	"ZADD",
	key,
	expires,
	member
)

return {
	was_session_online,
	was_user_online and 1 or 0
}
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
local member = ARGV[2]
local session_prefix = ARGV[3]

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

for _, current_member in ipairs(members) do
	if string.sub(
		current_member,
		1,
		string.len(session_prefix)
	) == session_prefix then
		remaining_session =
			remaining_session + 1
	end
end

local remaining_user =
	redis.call("ZCARD", key)

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
