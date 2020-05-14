//
//  Regular.h
//  Tinder
//
//  Created by Layer on 2020/5/2.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface Regular : NSObject

+(BOOL) isValidName:(NSString*)name; // 20位内字符
+(BOOL) isValidEmail:(NSString*)email; // 邮箱
+(BOOL) isValidPassword:(NSString*)password; // 密码

@end

NS_ASSUME_NONNULL_END
