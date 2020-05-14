//
//  UserDefaults.m
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "UserDefaults.h"

@implementation UserDefaults

+ (instancetype)standardUserDefaults {
    static dispatch_once_t onceToken;
    static UserDefaults* sharedInstance = nil;
    dispatch_once(&onceToken, ^{
        sharedInstance = [[self alloc] init];
    });
    return sharedInstance;
}

+ (void)putUserDefaults:(NSString*) key Value:(id)value {
    if (key != nil && value != nil) {
        NSUserDefaults* userDefaults = [NSUserDefaults standardUserDefaults];
        [userDefaults setObject:value forKey:key];
    }
}

+ (id)getUserDefaults:(NSString*)key {
    if (key != nil) {
        NSUserDefaults* userDefault = [NSUserDefaults standardUserDefaults];
        id obj = [userDefault objectForKey:key];
        return obj;
    }
    return nil;
}

@end
