//
//  LoginView.m
//  Tinder
//
//  Created by Layer on 2020/5/5.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "LoginView.h"

@implementation LoginView

- (instancetype)initWithFrame:(CGRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.emailTextField];
        [self addSubview:self.passwordTextFiled];
        [self addSubview:self.loginButton];
        [self addSubview:self.registerButton];
    }
    return self;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    _emailTextField.sd_layout.leftSpaceToView(self, 50).topSpaceToView(self, self.bounds.size.height / 2 - 100).rightSpaceToView(self, 50).heightIs(50);
    _passwordTextFiled.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_emailTextField, 10).rightSpaceToView(self, 50).heightIs(50);
    _loginButton.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_passwordTextFiled, 10).rightSpaceToView(self, 50).heightIs(50);
    _registerButton.sd_layout.leftSpaceToView(self, 50).rightSpaceToView(self, 50).bottomSpaceToView(self, 30).heightIs(50);
}

- (UITextField*)emailTextField {
    if (!_emailTextField) {
        _emailTextField = [[UITextField alloc] init];   // 构造
        _emailTextField.backgroundColor = [UIColor whiteColor]; // 背景
        _emailTextField.font = [UIFont systemFontOfSize:14]; // 字体
        _emailTextField.layer.cornerRadius = 20; // 圆角
        _emailTextField.layer.masksToBounds = YES;
        _emailTextField.placeholder = @"Enter email"; // 提示
        
        _emailTextField.leftView = [[UIView alloc]  initWithFrame:CGRectMake(0, 0, 20, 50)]; // 间隙
        _emailTextField.leftViewMode = UITextFieldViewModeAlways;
        _emailTextField.clearButtonMode = UITextFieldViewModeWhileEditing; // 删除键
        
        _emailTextField.tag = 1;
        _emailTextField.delegate = self;
    }
    return _emailTextField;
}

- (UITextField*)passwordTextFiled {
    if (!_passwordTextFiled) {
        _passwordTextFiled = [[UITextField alloc] init];
        _passwordTextFiled.backgroundColor = [UIColor whiteColor];
        _passwordTextFiled.font = [UIFont systemFontOfSize:14];
        _passwordTextFiled.layer.cornerRadius = 20;
        _passwordTextFiled.layer.masksToBounds = YES;
        _passwordTextFiled.placeholder = @"Enter password";
        
        _passwordTextFiled.leftView = [[UIView alloc] initWithFrame:CGRectMake(0, 0, 20, 50)];
        
        _passwordTextFiled.leftViewMode = UITextFieldViewModeAlways;
        _passwordTextFiled.clearButtonMode = UITextFieldViewModeWhileEditing;
        
        _passwordTextFiled.secureTextEntry = YES;
        
        _passwordTextFiled.tag = 2;
        _passwordTextFiled.delegate = self;
    }
    return _passwordTextFiled;
}

- (UIButton*)loginButton {
    if (!_loginButton) {
        _loginButton = [UIButton buttonWithType:UIButtonTypeSystem];
        _loginButton.backgroundColor = [UIColor colorWithRed:0.8 green:0 blue:0.3 alpha:1];
        [_loginButton setTitle:@"Log In" forState:UIControlStateNormal];
        [_loginButton setTintColor:[UIColor whiteColor]];
        _loginButton.titleLabel.font = [UIFont systemFontOfSize:15 weight:1];
        _loginButton.layer.cornerRadius = 20;
        _loginButton.layer.masksToBounds = YES;
        
        _loginButton.tag = 1;
        [_loginButton addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _loginButton;
}

-(UIButton*)registerButton {
    if (!_registerButton) {
        _registerButton = [UIButton buttonWithType:UIButtonTypeSystem];
        [_registerButton setTitleColor:[UIColor whiteColor] forState:UIControlStateNormal];
        [_registerButton setTitle:@"Back to Reginster" forState:UIControlStateNormal];
        _registerButton.titleLabel.font = [UIFont systemFontOfSize:15 weight:1];
        
        _registerButton.tag = 2;
        [_registerButton addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _registerButton;
}

-(void)textFieldDidBeginEditing:(UITextField *)textField {
    textField.textColor = [UIColor blackColor];
}

- (void)textFieldDidEndEditing:(UITextField *)textField {
    switch (textField.tag) {
        case 1: // NSLog(@"email");
            if (![Regular isValidEmail:textField.text]) {
                textField.textColor = [UIColor redColor];
            }
            break;
        case 2: // NSLog(@"password");
            if (![Regular isValidPassword:textField.text]) {
                textField.textColor = [UIColor redColor];
            }
            break;
        default:
            break;
    }
}

- (void)clickBtn:(UIButton*)button {
    if (button.tag == 1) {
        if (!([Regular isValidEmail:self.emailTextField.text] && [Regular isValidPassword:self.passwordTextFiled.text])) {
//            NSLog(@"缺失登录信息");
            [SVProgressHUD dismissWithDelay:1.5];
            [SVProgressHUD showInfoWithStatus:@"请检查信息!"];
            return;
        }
    }
    if ([_delegate respondsToSelector:@selector(clickButton:)]) {
        [_delegate clickButton:button.tag];
    }
}

@end
