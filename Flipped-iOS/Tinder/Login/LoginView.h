//
//  LoginView.h
//  Tinder
//
//  Created by Layer on 2020/5/5.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@protocol LoginDelegate <NSObject>

- (void)clickButton:(NSInteger)BtnTag;

@end

@interface LoginView : UIView <UITextFieldDelegate>

@property (strong, nonatomic) UITextField* emailTextField;  // 邮箱
@property (strong, nonatomic) UITextField* passwordTextFiled; // 密码
@property (strong, nonatomic) UIButton* loginButton; // 登录
@property (strong, nonatomic) UIButton* registerButton; // 注册

@property (weak, nonatomic) id<LoginDelegate> delegate; // 代理

@end

NS_ASSUME_NONNULL_END
