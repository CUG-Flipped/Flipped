//
//  UserDefaults.h
//  Tinder
//
//  Created by Layer on 2020/5/6.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface UserDefaults : NSObject

+ (instancetype)standardUserDefaults;
+ (void)putUserDefaults:(NSString*) key Value:(id)value;
+ (id)getUserDefaults:(NSString*)key;

@end

NS_ASSUME_NONNULL_END
