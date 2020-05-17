//
//  HttpData.h
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@class UIImage;
@interface HttpData : NSObject

// 注册
+ (void)reginsterWithUserType:(NSString*)user_type withName:(NSString*)name withEmail:(NSString*)email withPhoto:(NSString*)photo withPassword:(NSString*)password seccess:(void(^)(id json))success failure:(void (^)(NSError* err))failure;

// 登录
+ (void)loginWithUserType:(NSString*)user_type withEmail:(NSString*)email withPassword:(NSString*)password success:(void(^)(id json))success failure:(void(^)(NSError* err))failure;

// 文件上传，获取路径url（图片）
+ (void)requestImgUploadFile:(UIImage *)file success:(void(^)(id json))success failure:(void(^)(NSError *err))failure;

// 详细信息
+ (void)requestHomeInfoToken:(NSString*)token user_id:(NSString*)user_id success:(void(^)(id json))success failure:(void(^)(NSError* err))failure;

// 首页列表
+ (void)requestHomeListToken:(NSString*)token page:(NSInteger)page pageSize:(NSInteger)pageSize success:(void(^)(id json))success failure:(void(^)(NSError* err))failure;

// 首页收藏
+ (void)requestHomeCollectToken:(NSString*)token user_id:(NSString*)user_id type:(NSString*)type success:(void(^)(id json))success failure:(void(^)(NSError* err))failure;



@end

NS_ASSUME_NONNULL_END
