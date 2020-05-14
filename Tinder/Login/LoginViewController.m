//
//  LoginViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/2.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "LoginViewController.h"

@interface LoginViewController () <LoginDelegate>

@property (strong, nonatomic) LoginView* loginView;

@end

@implementation LoginViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    [self createLoginView];
   
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillShowNotification object:nil];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillHideNotification object:nil];
}

# pragma mark - 创建LoginView
- (void)createLoginView {
    self.loginView = [[LoginView alloc] initWithFrame:self.view.bounds];
    [self.view addSubview: self.loginView];
    self.loginView.delegate = self;
    
    // 渐变
    CAGradientLayer* gradient = [CAGradientLayer layer];
    gradient.frame = self.view.bounds;
    NSMutableArray* graColor = [[NSMutableArray alloc] initWithObjects:(id)[UIColor colorWithRed:250 / 255.0 green:94 / 255.0 blue:99 / 255.0 alpha:1].CGColor, (id)[UIColor colorWithRed:227 / 255.0 green:40 / 255.0 blue:117 / 255.0 alpha:1].CGColor, nil];
    gradient.colors = graColor;
    
    gradient.startPoint = CGPointMake(0, 0);
    gradient.endPoint = CGPointMake(0, 1);
    [self.view.layer insertSublayer:gradient atIndex:0];
}

# pragma mark - 点击空白 退出编辑
- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [self.loginView.emailTextField resignFirstResponder];
    [self.loginView.passwordTextFiled resignFirstResponder];
}

# pragma mark - 键盘
- (void)keyboardShowAndHide:(NSNotification*)notification {
    NSDictionary* info = [notification userInfo];
    CGRect beginRect = [[info objectForKey:UIKeyboardFrameBeginUserInfoKey] CGRectValue];
    CGRect endRect = [[info objectForKey:UIKeyboardFrameEndUserInfoKey] CGRectValue];
    CGFloat change = (endRect.origin.y - beginRect.origin.y) / 2;
    [UIView animateWithDuration:0.25 animations:^{
        [self.view setFrame:CGRectMake(self.view.frame.origin.x, self.view.frame.origin.y + change, self.view.frame.size.width, self.view.frame.size.height)];
    }];
}

#pragma mark - 按键响应
- (void)clickButton:(NSInteger)BtnTag {
    switch (BtnTag) {
        case 1: // NSLog(@"1");
            [self clickLogin];
            break;
        case 2: // NSLog(@"2");
            [self clickGotoReginster];
            break;
        default:
            break;
    }
}
- (void)clickLogin {
    NSLog(@"登录");
    [HttpData loginWithUserType:@"1" withEmail:self.loginView.emailTextField.text withPassword:self.loginView.passwordTextFiled.text success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"]integerValue] == 200) {
                [UserDefaults putUserDefaults:@"email" Value:self.loginView.emailTextField.text];
                [UserDefaults putUserDefaults:@"password" Value:self.loginView.passwordTextFiled.text];
                [UserDefaults putUserDefaults:@"isLogin" Value:@"YES"];
                [UserDefaults putUserDefaults:@"token" Value:json[@"data"][@"token"]];
                [UserDefaults putUserDefaults:@"user_id" Value:json[@"data"][@"id"]];
                
//                NSLog(@"成功");
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showSuccessWithStatus:@"登录成功!"];
                
                HomeViewController* homeVC = [[HomeViewController alloc] init];
                [self.navigationController pushViewController:homeVC animated:YES];
            }
            else {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:json[@"message"]];
            }
        }
    } failure:^(NSError * _Nonnull err) {
//        NSLog(@"失败");
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showErrorWithStatus:@"登录失败!"];
    }];
}

- (void)clickGotoReginster {
//    [self.navigationController popViewControllerAnimated:YES];
    [self.navigationController pushViewController:[[ViewController alloc] init] animated:YES];
}


@end
