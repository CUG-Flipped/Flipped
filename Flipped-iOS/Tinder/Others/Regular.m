//
//  Regular.m
//  Tinder
//
//  Created by Layer on 2020/5/2.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "Regular.h"

@implementation Regular

+(BOOL) isValidName:(NSString*)name {
    NSString *str = @"^[a-zA-z\u4E00-\u9FA5]{5,20}";
    NSPredicate* prediate = [NSPredicate predicateWithFormat:@"SELF MATCHES %@", str];
    return [prediate evaluateWithObject:name];
}
+(BOOL) isValidEmail:(NSString*)email {
    NSString *str = @"[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}";
    NSPredicate* prediate = [NSPredicate predicateWithFormat:@"SELF MATCHES %@", str];
    return [prediate evaluateWithObject:email];
}
+(BOOL) isValidPassword:(NSString*)password {
    if (password.length >= 6 && password.length <= 20) return YES;
    else return NO;
}

@end
