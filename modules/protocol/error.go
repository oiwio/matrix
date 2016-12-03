package protocol

const (
	ERROR_BASE                      = 0
	ERROR_INTERNAL_ERROR            = ERROR_BASE - 1
	ERROR_DEVICE_ALREADY_REGISTRIED = ERROR_BASE - 2
	ERROR_USERNAME_ALREADY_INUSE    = ERROR_BASE - 3  // -2 用户名已被注册
	ERROR_PHONE_ALREADY_REGISTRIED  = ERROR_BASE - 4  // -3 手机号已经被注册
	ERROR_EMAIL_ALREADY_REGISTRIED  = ERROR_BASE - 5  // -4 邮箱已经被注册
	ERROR_CANNOT_REGISTRY           = ERROR_BASE - 6  // -5 注册发生错误
	ERROR_PASSWORD_NOT_MATCH        = ERROR_BASE - 7  // -6 用户名密码不匹配
	ERROR_NEED_SIGNIN               = ERROR_BASE - 8  // -7 需要登陆
	ERROR_USER_ALREADY_SINGNED_IN   = ERROR_BASE - 9  // -8 用户已经登录
	ERROR_WRONG_VERIFY_CODE         = ERROR_BASE - 10 // -9 错误的验证码
	ERROR_INVALID_REQUEST           = ERROR_BASE - 11 // -10 错误的请求
	ERROR_INVALID_PHONE_NUMBER      = ERROR_BASE - 12 // -11 无效的手机号
	ERROR_INVALID_EMAIL_NUMBER      = ERROR_BASE - 13 // -12 无效的邮箱
	ERROR_INVALID_USER              = ERROR_BASE - 14 // -13 用户不存在
	ERROR_INVALID_CODE              = ERROR_BASE - 15 // -14 失效的验证码
	ERROR_INVALID_ACCOUNT           = ERROR_BASE - 16 // -15 无效的账号
	ERROR_INVALID_GENDER            = ERROR_BASE - 17 // -16 无效的性别
	ERROR_INVALID_DEVICETOKEN       = ERROR_BASE - 18 // -17 无效的设备Token
	ERROR_NOT_WHISPER               = ERROR_BASE - 19 //对方设置不能进行私聊
)
