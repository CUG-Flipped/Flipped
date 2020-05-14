//
//  HomeView.m
//  Tinder
//
//  Created by Layer on 2020/4/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "HomeView.h"

@implementation HomeView

# pragma mark - 初始化头像
- (UIImageView*)headImage {
    if (!_headImage) {
        _headImage = [[UIImageView alloc] init]; // 构造
        _headImage.backgroundColor = [UIColor whiteColor]; // 背景
        _headImage.layer.cornerRadius = 10; // 圆角
        _headImage.layer.masksToBounds = YES;
        
        // 允许交互
        _headImage.userInteractionEnabled = YES;
        
        // 点击手势
        UITapGestureRecognizer* tapGes = [[UITapGestureRecognizer alloc] initWithTarget:self action:@selector(clickImg)];
        [_headImage addGestureRecognizer:tapGes];
    }
    return _headImage;
}

# pragma mark - 初始化头像提示
- (UILabel*)headLabel {
    if (!_headLabel) {
        _headLabel = [[UILabel alloc] init]; // 构造
        _headLabel.text = @"Select Photo"; // 文字
        _headLabel.font = [UIFont systemFontOfSize:23 weight:1]; // 字体、加粗
        _headLabel.textAlignment = NSTextAlignmentCenter; // 居中
    }
    return _headLabel;
}

# pragma mark - 初始化名字输入
- (UITextField*)nameTextFiled {
    if (!_nameTextFiled) {
        _nameTextFiled = [[UITextField alloc] init]; // 构造
        _nameTextFiled.backgroundColor = [UIColor whiteColor]; // 背景
        _nameTextFiled.placeholder = @"Enter full name"; // 提示
        _nameTextFiled.layer.cornerRadius = 20; // 圆角
        _nameTextFiled.layer.masksToBounds = YES;
        _nameTextFiled.font = [UIFont systemFontOfSize:14]; // 字体
        _nameTextFiled.leftView = [[UIView alloc] initWithFrame:CGRectMake(0, 0, 20, 50)]; // 左边间隙
        _nameTextFiled.leftViewMode = UITextFieldViewModeAlways; // 左边间隙一直显示
        _nameTextFiled.clearButtonMode = UITextFieldViewModeWhileEditing; // 删除键
        
        _nameTextFiled.tag = 1;
        _nameTextFiled.delegate = self;
    }
    return _nameTextFiled;
}

# pragma mark - 初始化邮箱
- (UITextField*)emailTextFiled {
    if (!_emailTextFiled) {
        _emailTextFiled = [[UITextField alloc] init]; // 构造
        _emailTextFiled.backgroundColor = [UIColor whiteColor]; // 背景
        _emailTextFiled.placeholder = @"Enter email"; // 提示
        _emailTextFiled.layer.cornerRadius = 20; // 圆角
        _emailTextFiled.layer.masksToBounds = YES;
        _emailTextFiled.font = [UIFont systemFontOfSize:14]; // 字体
        _emailTextFiled.leftView = [[UIView alloc] initWithFrame:CGRectMake(0, 0, 20, 50)]; // 左边间隙
        _emailTextFiled.leftViewMode = UITextFieldViewModeAlways; // 左边间隙一直显示
        _emailTextFiled.clearButtonMode = UITextFieldViewModeWhileEditing; // 删除键
        
        _emailTextFiled.tag = 2;
        _emailTextFiled.delegate = self;
    }
    return _emailTextFiled;
}

# pragma mark - 初始化密码
- (UITextField*)passwordTextFiled {
    if (!_passwordTextFiled) {
        _passwordTextFiled = [[UITextField alloc] init]; // 构造
        _passwordTextFiled.backgroundColor = [UIColor whiteColor]; // 背景
        _passwordTextFiled.placeholder = @"Enter password"; // 提示
        _passwordTextFiled.layer.cornerRadius = 20; // 圆角
        _passwordTextFiled.layer.masksToBounds = YES;
        _passwordTextFiled.font = [UIFont systemFontOfSize:14]; // 字体
        _passwordTextFiled.leftView = [[UIView alloc] initWithFrame:CGRectMake(0, 0, 20, 50)]; // 左边间隙
        _passwordTextFiled.leftViewMode = UITextFieldViewModeAlways; // 左边间隙一直显示
        _passwordTextFiled.clearButtonMode = UITextFieldViewModeWhileEditing; // 删除键
        _passwordTextFiled.secureTextEntry = YES; // 密码隐藏
        
        _passwordTextFiled.tag = 3;
        _passwordTextFiled.delegate = self;
    }
    return _passwordTextFiled;
}

# pragma mark - 初始化注册
- (UIButton*)registerButton {
    if (!_registerButton) {
        _registerButton = [UIButton buttonWithType:UIButtonTypeSystem];
        _registerButton.backgroundColor = [UIColor colorWithRed:0.8 green:0 blue:0.3 alpha:1]; // 背景
        [_registerButton setTitle:@"Register" forState:UIControlStateNormal]; // 正常状态
        [_registerButton setTitleColor:[UIColor whiteColor] forState:UIControlStateNormal]; // 颜色
        _registerButton.titleLabel.font = [UIFont systemFontOfSize:15 weight:1]; // 字体、加粗
        _registerButton.layer.cornerRadius = 20; // 圆角
        _registerButton.layer.masksToBounds = YES;
        
        _registerButton.tag = 1;
        [_registerButton addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _registerButton;
}

# pragma mark - 初始化登录
- (UIButton*)loginButton {
    if (!_loginButton) {
        _loginButton = [UIButton buttonWithType:UIButtonTypeSystem];
        [_loginButton setTitle:@"Go to Login" forState:UIControlStateNormal]; // 正常状态
        [_loginButton setTitleColor:[UIColor whiteColor] forState:UIControlStateNormal]; // 颜色
        _loginButton.titleLabel.font = [UIFont systemFontOfSize:15 weight:1]; // 字体、加粗
        
        _loginButton.tag = 2;
        [_loginButton addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _loginButton;
}

# pragma mark - 构造函数
- (instancetype)initWithFrame:(CGRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.headImage];
        [self.headImage addSubview:self.headLabel];
        [self addSubview:self.nameTextFiled];
        [self addSubview:self.emailTextFiled];
        [self addSubview:self.passwordTextFiled];
        [self addSubview:self.registerButton];
        [self addSubview:self.loginButton];
    }
    return self;
}

# pragma mark - 位置
- (void)layoutSubviews {
    [super layoutSubviews];
    _headImage.sd_layout.leftSpaceToView(self, 50).topSpaceToView(self, 130).rightSpaceToView(self, 50).heightIs(self.bounds.size.width - 100);
    _headLabel.sd_layout.leftSpaceToView(_headImage, 0).topSpaceToView(_headImage, 0).rightSpaceToView(_headImage, 0).bottomSpaceToView(_headImage, 0);
    _nameTextFiled.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_headImage, 20).rightSpaceToView(self, 50).heightIs(50);
    _emailTextFiled.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_nameTextFiled, 10).rightSpaceToView(self, 50).heightIs(50);
    _passwordTextFiled.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_emailTextFiled, 10).rightSpaceToView(self, 50).heightIs(50);
    _registerButton.sd_layout.leftSpaceToView(self, 50).topSpaceToView(_passwordTextFiled, 10).rightSpaceToView(self, 50).heightIs(50);
    _loginButton.sd_layout.leftSpaceToView(self, 50).rightSpaceToView(self, 50).bottomSpaceToView(self, 30).heightIs(50);
}

# pragma mark - 点击图片
- (void)clickImg {
    if ([_delegate respondsToSelector:@selector(selectImage)]) {
        [_delegate selectImage];
    }
}

# pragma mark - 点击空白
- (void)textFieldDidBeginEditing:(UITextField *)textField {
    textField.textColor = [UIColor blackColor];
}

# pragma mark - 输入完成后
- (void)textFieldDidEndEditing:(UITextField *)textField {
    switch(textField.tag) {
        case 1: // 名字
            if (![Regular isValidName:textField.text]) textField.textColor = [UIColor redColor];
            else textField.textColor = [UIColor blackColor];
            break;
        case 2: // 邮箱
            if (![Regular isValidEmail:textField.text]) textField.textColor = [UIColor redColor];
            else textField.textColor = [UIColor blackColor];
            break;
        case 3: // 密码
            if (![Regular isValidPassword:textField.text]) textField.textColor = [UIColor redColor];
            else textField.textColor = [UIColor blackColor];
            break;
        default:
            break;
    }
}

# pragma mark - 按键响应
-(void)clickBtn:(UIButton*)button {
    if (button.tag == 1) {
        if (!([Regular isValidName:self.nameTextFiled.text] && [Regular isValidEmail:self.emailTextFiled.text] && [Regular isValidPassword:self.passwordTextFiled.text])) {
//            NSLog(@"缺失注册信息");
            [SVProgressHUD dismissWithDelay:1];
            [SVProgressHUD showInfoWithStatus:@"请检查信息!"];
//            [SVProgressHUD showImage:[UIImage imageNamed:@""] status:@"请检查注册信息！"];
            return;
        }
    }
    if ([_delegate respondsToSelector:@selector(clickButton:)]) {
        [_delegate clickButton:button.tag];
    }
}

@end
