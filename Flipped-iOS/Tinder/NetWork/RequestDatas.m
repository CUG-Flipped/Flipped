//
//  RequestDatas.m
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "RequestDatas.h"

static RequestDatas* request = nil;

@implementation RequestDatas

# pragma mark - 单例
+ (instancetype)shareRequestDatas {
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        request = [[RequestDatas alloc] init];
    });
    return request;
}

+ (instancetype)allocWithZone:(struct _NSZone *)zone {
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        request = [super allocWithZone:zone];
    });
    return request;
}

# pragma mark - 网络请求管理者
+ (AFHTTPSessionManager*)Manager {
    AFHTTPSessionManager* manager = [[AFHTTPSessionManager alloc] initWithBaseURL:[NSURL URLWithString:@"http://101.132.114.122:8099"]]; // 服务器前缀
    
    AFJSONResponseSerializer* response = [AFJSONResponseSerializer serializer]; // 过滤空数据
    response.removesKeysWithNullValues = YES;
    manager.responseSerializer = response;
    
//    manager.requestSerializer = [AFJSONRequestSerializer serializer]; // 请求数据格式
    [manager.requestSerializer setValue:@"application/json" forHTTPHeaderField:@"Content－Type"];
    
    manager.requestSerializer.timeoutInterval = 10; // 请求时间
    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObjects:@"text/html", @"text/plain", @"application/json", @"image/jpeg",@"image/png", nil]; // 接收类型

    return manager;
}

# pragma mark - 响应
+ (void)requestWithMethond:(RequestMethondType)type withUrl:(NSString*)url  withIsPullDown:(BOOL)isPullDown  withParamter:(NSDictionary*)paramter seccess:(void (^)(id resoponse))success failure:(void (^)(NSError* err))failure {
    
    AFHTTPSessionManager *manager = [self Manager];
    
    switch (type) {
        case RequestMethondTypeGet: {
            [manager GET:url parameters:paramter progress:^(NSProgress * _Nonnull downloadProgress) {
                // 进度
            } success:^(NSURLSessionDataTask * _Nonnull task, id  _Nullable responseObject) {
                if  (success) {
                    success(responseObject);
                }
            } failure:^(NSURLSessionDataTask * _Nullable task, NSError * _Nonnull error) {
                if (failure) {
                    failure(error);
                }
            }];
            
            break;
        }
            
        case RequestMethondTypePost: {
            [manager POST:url parameters:paramter progress:^(NSProgress * _Nonnull uploadProgress) {
                // 进度
            } success:^(NSURLSessionDataTask * _Nonnull task, id  _Nullable responseObject) {
                if (success) {
                    success(responseObject);
                }
            } failure:^(NSURLSessionDataTask * _Nullable task, NSError * _Nonnull error) {
                if (failure) {
                    failure(error);
                }
            }];
            break;
        }
            
        case RequestMethondTypeImage: {
            [manager POST:url parameters:paramter constructingBodyWithBlock:^(id<AFMultipartFormData>  _Nonnull formData) {
                
                NSDate* data = paramter[@""];
                
                if (!data) {
                    data = [NSData data];
                }
                
                [formData appendPartWithFileData:data name:@"headPortrait" fileName:@"headPortrait" mimeType:@"image/jpeg"];
                
            } progress:^(NSProgress * _Nonnull uploadProgress) {
                // 进度
            } success:^(NSURLSessionDataTask * _Nonnull task, id  _Nullable responseObject) {
                if (success) {
                    success(responseObject);
                }
            } failure:^(NSURLSessionDataTask * _Nullable task, NSError * _Nonnull error) {
                if (failure) {
                    failure(error);
                }
            }];
            break;
        }
           
        default:
            break;
    }
}




@end
