# Matchee All APIs - Unified Postman Collection

## 📁 File Tổng Hợp

**`Matchee_All_APIs.postman_collection.json`** - File duy nhất chứa tất cả:
- ✅ Authentication APIs
- ✅ User Management APIs  
- ✅ Player Profile APIs
- ✅ Player Suggestions APIs
- ✅ Complete Test Flows
- ✅ Environment Variables
- ✅ Auto-save Scripts
- ✅ New Response Format Support

## 🚀 Cách Sử Dụng

### Bước 1: Import Collection
1. Mở Postman
2. Click **Import** button
3. Chọn file `Matchee_All_APIs.postman_collection.json`
4. Click **Import**

### Bước 2: Chọn Environment
1. Click dropdown environment ở góc trên bên phải
2. Chọn **Matchee All APIs Environment** (tự động tạo từ collection)

## 🎯 Cấu Trúc Collection

### 1. Health Check
- **GET /health** - Kiểm tra trạng thái server

### 2. Authentication
- **POST /api/v1/auth/register** - Đăng ký user mới
- **POST /api/v1/auth/login** - Đăng nhập
- **POST /api/v1/auth/refresh** - Làm mới token
- **POST /api/v1/auth/logout** - Đăng xuất

### 3. User Management
- **GET /api/v1/users/me** - Lấy thông tin user hiện tại
- **PUT /api/v1/users/me** - Cập nhật profile
- **POST /api/v1/users/change-password** - Đổi mật khẩu

### 4. Player Profile Management
- **POST /api/v1/player/profile** - Tạo/cập nhật player profile
- **GET /api/v1/player/profile** - Lấy player profile hiện tại
- **GET /api/v1/player/profile/:id** - Lấy player profile theo ID

### 5. Player Suggestions
- **GET /api/v1/player/suggestions** - Lấy gợi ý đối thủ
  - Query params: `level`, `latitude`, `longitude`, `radius`, `limit`

### 6. Complete Test Flows
- **Authentication Flow** - Flow đầy đủ cho auth
- **Player Profile Flow** - Flow đầy đủ cho player profile

## 🔧 Environment Variables

Collection tự động tạo environment với các variables:

### Basic Variables
- `base_url`: http://localhost:8080
- `access_token`: (auto-filled)
- `refresh_token`: (auto-filled)
- `user_id`: (auto-filled)
- `player_profile_id`: (auto-filled)

### User Variables
- `user_full_name`: Test User
- `user_phone`: 0123456789
- `user_email`: test@example.com
- `user_password`: password123

### Update Variables
- `updated_full_name`: Updated Test User
- `updated_email`: updated@example.com
- `updated_phone`: 0987654321
- `new_password`: newpassword456

### Player Variables
- `player_level`: good
- `player_gender`: male
- `player_location`: Ho Chi Minh City
- `player_latitude`: 10.8231
- `player_longitude`: 106.6297
- `player_bio`: Passionate tennis player

### Suggestion Variables
- `suggestion_level`: good
- `suggestion_latitude`: 10.8231
- `suggestion_longitude`: 106.6297
- `suggestion_radius`: 10
- `suggestion_limit`: 5

## 🧪 Test Features

### Auto-save Scripts
- ✅ Tự động lưu tokens từ login/register
- ✅ Tự động lưu user_id và player_profile_id
- ✅ Hỗ trợ cả response format cũ và mới
- ✅ Debug logging cho troubleshooting

### New Response Format Support
```json
{
  "status": "success" | "fail",
  "code": 200 | 400 | 401 | 404 | 409 | 500,
  "message": "string",
  "data": any | null
}
```

### Test Scripts Features
- ✅ Validate new response format
- ✅ Auto-extract tokens from data field
- ✅ Fallback support for old format
- ✅ Console logging for debugging
- ✅ Success/error detection

## 📋 Test Scenarios

### Scenario 1: Authentication Flow
```
1. Health Check
2. Register User
3. Login User
4. Get Current User
5. Update Profile
6. Change Password
7. Logout
```

### Scenario 2: Player Profile Flow
```
1. Health Check
2. Login User
3. Create Player Profile
4. Get Player Profile
5. Get Player Suggestions
```

## 🎯 Quick Start

### Option 1: Use Complete Test Flows
1. Mở **Complete Test Flows**
2. Chọn **Authentication Flow** hoặc **Player Profile Flow**
3. Click **Run** để chạy toàn bộ flow

### Option 2: Test Individual APIs
1. Chọn API cần test
2. Đảm bảo đã có access_token (từ login)
3. Click **Send**

## 🔍 Debugging

### Check Environment Variables
1. Click **Environments** tab
2. Select **Matchee All APIs Environment**
3. Verify all variables are set correctly

### Check Response Logs
1. Open **Console** (View → Show Postman Console)
2. Check logs for debugging info
3. Verify token extraction

### Response Format Validation
- ✅ New format: `{status, code, message, data}`
- ✅ Old format: Direct response object
- ✅ Auto-detection and handling

## 💡 Tips

1. **Use Test Flows**: Sử dụng Complete Test Flows để test toàn bộ hệ thống
2. **Check Console**: Luôn check console logs để debug
3. **Environment Variables**: Tất cả variables đã được setup sẵn
4. **Auto-save**: Tokens và IDs được tự động lưu
5. **Response Format**: Hỗ trợ cả format cũ và mới

## 🚨 Common Issues

### Token Not Saved
- Check if login response contains `accessToken` in `data` field
- Verify response format is correct
- Check console logs for errors

### 401 Unauthorized
- Verify token is valid and not expired
- Check if token is properly set in Authorization header
- Try refreshing token

### Empty Suggestions
- Check if there are other players in database
- Verify location parameters are correct
- Check if user has player profile

## 📝 Migration Notes

- ✅ Tất cả APIs đã được cập nhật với response format mới
- ✅ Test scripts đã được cập nhật để handle format mới
- ✅ Environment variables đã được tối ưu
- ✅ Complete test flows đã được tạo
- ✅ Documentation đã được cập nhật

## 🎉 Benefits

1. **Single File**: Chỉ cần 1 file để quản lý tất cả
2. **Complete Coverage**: Bao gồm tất cả APIs và test flows
3. **Auto-save**: Tự động lưu tokens và IDs
4. **Format Support**: Hỗ trợ cả response format cũ và mới
5. **Easy Testing**: Test flows có sẵn để test toàn bộ hệ thống
6. **Debugging**: Console logs và validation scripts
7. **Maintenance**: Dễ dàng cập nhật và maintain

Từ giờ chỉ cần cập nhật file `Matchee_All_APIs.postman_collection.json` duy nhất! 🎯📊
