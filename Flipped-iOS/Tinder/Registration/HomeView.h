//
//  HomeView.h
//  Tinder
//
//  Created by Layer on 2020/4/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

// 代理
@protocol RegisterDelegate <NSObject>

- (void)selectImage;
- (void)clickButton:(NSInteger)btnTag;

@end

@interface HomeView : UIView <UITextFieldDelegate>

@property (strong, nonatomic) UIImageView* headImage; // 头像
@property (strong, nonatomic) UILabel* headLabel; // 头像提示
@property (strong, nonatomic) UITextField* nameTextFiled; // 名字
@property (strong, nonatomic) UITextField* emailTextFiled; // 邮箱
@property (strong, nonatomic) UITextField* passwordTextFiled; // 密码
@property (strong, nonatomic) UIButton* registerButton; // 注册
@property (strong, nonatomic) UIButton* loginButton; // 登录

@property (weak, nonatomic) id<RegisterDelegate> delegate; // 代理

@end

NS_ASSUME_NONNULL_END
