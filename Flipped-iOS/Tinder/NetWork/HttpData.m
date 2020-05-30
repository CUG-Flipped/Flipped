//
//  HttpData.m
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "HttpData.h"

@implementation HttpData

+ (void)reginsterWithUserType:(NSString*)user_type withName:(NSString*)name withEmail:(NSString*)email withPhoto:(NSString*)photo withPassword:(NSString*)password seccess:(void(^)(id json))success failure:(void (^)(NSError* err))failure {
    
    NSDictionary *dic = @{
        @"user_type":user_type,
        @"name":name,
        @"email":email,
        @"photo":photo,
        @"password":password
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypePost withUrl:@"/user/register" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
    
}

+ (void)loginWithUserType:(NSString*)user_type withEmail:(NSString*)email withPassword:(NSString*)password success:(void(^)(id json))success failure:(void(^)(NSError*err))failure {
    
    NSDictionary* dic = @ {
        @"user_type":user_type,
        @"email":email,
        @"password":password
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypePost withUrl:@"/user/login" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
}

/** 文件上传，获取转换后的路径url（图片） */
+(void)requestImgUploadFile:(UIImage *)file success:(void(^)(id json))success failure:(void(^)(NSError *err))failure
{
    NSData *imageData;
    NSString *imageFormat;
    if (UIImagePNGRepresentation(file) != nil) {
        imageFormat = @"Content-Type: image/png \r\n";
        imageData = UIImagePNGRepresentation(file);
        
    }else{
        imageFormat = @"Content-Type: image/jpeg \r\n";
        imageData = UIImageJPEGRepresentation(file, 0.7);
    }
    NSURL *url = [NSURL URLWithString:@"http://101.132.114.122:8099/upload"];
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:url];
    request.HTTPMethod = @"POST";
    //设置请求实体
    NSMutableData *body = [NSMutableData data];
    
    //文件参数
    [body appendData:[self getDataWithString:@"--BOUNDARY\r\n" ]];
    NSString *disposition = [NSString stringWithFormat:@"Content-Disposition: form-data; name=\"file\"; filename=\"file.jpg\"\r\n"];
    [body appendData:[self getDataWithString:disposition ]];
    [body appendData:[self getDataWithString:imageFormat]];
    [body appendData:[self getDataWithString:@"\r\n"]];
    [body appendData:imageData];
    [body appendData:[self getDataWithString:@"\r\n"]];
    //普通参数
    [body appendData:[self getDataWithString:@"--BOUNDARY\r\n" ]];
    //上传参数需要key： （相应参数，在这里是file）
    NSString *dispositions = [NSString stringWithFormat:@"Content-Disposition: form-data; name=\"%@\"\r\n",@"key"];
    [body appendData:[self getDataWithString:dispositions ]];
    [body appendData:[self getDataWithString:@"\r\n"]];
    [body appendData:[self getDataWithString:@"file"]];
    [body appendData:[self getDataWithString:@"\r\n"]];
    
    //参数结束
    [body appendData:[self getDataWithString:@"--BOUNDARY--\r\n"]];
    request.HTTPBody = body;
    //设置请求体长度
    NSInteger length = [body length];
    [request setValue:[NSString stringWithFormat:@"%ld",length] forHTTPHeaderField:@"Content-Length"];
    //设置 POST请求文件上传
    [request setValue:@"multipart/form-data; boundary=BOUNDARY" forHTTPHeaderField:@"Content-Type"];
    
    //运用AFN实现照片上传
    AFHTTPSessionManager *manager = [self manager];
    
    NSDictionary *dict = @{@"file":file};
    
    [manager POST:@"http://101.132.114.122:8099/upload" parameters:dict constructingBodyWithBlock:^(id<AFMultipartFormData>  _Nonnull formData) {
        
        [formData appendPartWithFileData:imageData name:@"file" fileName:[NSString stringWithFormat:@"file.jpg"] mimeType:@"image/jpeg"];
        
    } progress:^(NSProgress * _Nonnull uploadProgress) {
        
    } success:^(NSURLSessionDataTask * _Nonnull task, id  _Nullable responseObject) {
        
        if(success)
        {
            success(responseObject);
        }
        
    } failure:^(NSURLSessionDataTask * _Nullable task, NSError * _Nonnull error) {
        
        if(failure)
        {
            failure(error);
        }
        
    }];
}

+(NSData *)getDataWithString:(NSString *)string{
    
    NSData *data = [string dataUsingEncoding:NSUTF8StringEncoding];
    return data;
    
}

+(AFHTTPSessionManager*) manager
{
    NSURL *mybaseURL = [NSURL URLWithString:baseUrl];
    // 1.创建请求管理对象
    AFHTTPSessionManager *manager = [[AFHTTPSessionManager alloc] initWithBaseURL:mybaseURL];
    // 2.过滤NSNull参数
    ((AFJSONResponseSerializer *)manager.responseSerializer).removesKeysWithNullValues = YES;
    
    // 3.请求的数据和接收的数据格式
    manager.responseSerializer =[AFJSONResponseSerializer serializer];
    [manager.requestSerializer setValue:@"multipart/form-data" forHTTPHeaderField:@"Content－Type"];
//    manager.requestSerializer  = [AFHTTPRequestSerializer serializer];
    // 4.设置请求成功后的接受内容的类型
    manager.responseSerializer.acceptableContentTypes = [NSSet setWithObjects:@"text/html", @"text/plain",@"application/json", @"image/jpeg",@"image/png",nil];
    manager.requestSerializer.timeoutInterval = 10;
    
    return manager;
}

// 详细信息
+ (void)requestHomeInfoToken:(NSString*)token user_id:(NSString*)user_id success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
    
    NSDictionary* dic = @{
        @"token" : token,
        @"user_id" : user_id
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypePost withUrl:@"/home/user_info" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
}

// 首页列表
+ (void)requestHomeListToken:(NSString*)token page:(NSInteger)page pageSize:(NSInteger)pageSize success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
    
    NSDictionary* dic = @{
        @"token" : token,
        @"page" : @(page),
        @"pageSize" : @(pageSize)
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypeGet withUrl:@"/user/public_user_list" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
}


// 首页收藏
+ (void)requestHomeCollectToken:(NSString*)token user_id:(NSString*)user_id type:(NSString*)type success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
    
    NSDictionary* dic = @{
        @"token" : token,
        @"user_id" : user_id,
        @"type" : type
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypePost withUrl:@"/user/operation" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
}

// 个人信息获取
+ (void)requestInfoToken:(NSString *)token success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
    
    NSDictionary* dic = @{
        @"token" : token
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypeGet withUrl:@"/user/personal" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
}

// 上传信息
+ (void)requestSaveOneSelfDataToken:(NSString *)token photo:(NSString *)photo photo1:(NSString *)photo1 photo2:(NSString *)photo2 photo3:(NSString *)photo3 user_name:(NSString *)user_name age:(NSString *)age work:(NSString *)work bio:(NSString *)bio age_min:(NSString *)age_min age_max:(NSString *)age_max success:(void(^)(id json))success failure:(void(^)(NSError *err))failure
{
    NSDictionary* parmas = @{
                             @"token":token,
                             @"photo":photo,
                             @"photo1":photo1,
                             @"photo2":photo2,
                             @"photo3":photo3,
                             @"user_name":user_name,
                             @"age":age,
                             @"work":work,
                             @"bio":bio,
                             @"age_min":age_min,
                             @"age_max":age_max
                             };
    
    [RequestDatas requestWithMethond:RequestMethondTypePost withUrl:@"/user/user_save" withIsPullDown:YES withParamter:parmas seccess:success failure:failure];
}

// 好友列表
+ (void)requestFriendDataToken:(NSString *)token success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
    
    NSDictionary* dic = @{
        @"token" : token
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypeGet withUrl:@"/message/message_list" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
    
}

+ (void)requestCollectFriendDataToken:(NSString *)token success:(void(^)(id json))success failure:(void(^)(NSError* err))failure {
   
    NSDictionary* dic = @{
        @"token" : token
    };
    
    [RequestDatas requestWithMethond:RequestMethondTypeGet withUrl:@"/message/collect_list" withIsPullDown:YES withParamter:dic seccess:success failure:failure];
    
}

@end

