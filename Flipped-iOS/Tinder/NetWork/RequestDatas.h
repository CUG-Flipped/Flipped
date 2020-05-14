//
//  RequestDatas.h
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <Foundation/Foundation.h>

typedef NS_ENUM(NSInteger, RequestMethondType) {
    RequestMethondTypeGet = 1,
    RequestMethondTypePost = 2,
    RequestMethondTypeImage = 3,
};

NS_ASSUME_NONNULL_BEGIN

@class AFHTTPSessionManager;
@interface RequestDatas : NSObject

// 单例
+ (instancetype)shareRequestDatas;

+ (void)requestWithMethond:(RequestMethondType)type withUrl:(NSString*)url withIsPullDown:(BOOL)isPullDown withParamter:(NSDictionary*)paramter seccess:(void (^)(id resoponse))success failure:(void (^)(NSError* err))failure;

+ (AFHTTPSessionManager*)Manager;

@end

NS_ASSUME_NONNULL_END
