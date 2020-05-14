//
//  ViewController.m
//  Tinder
//
//  Created by Layer on 2020/4/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "ViewController.h"

@interface ViewController () <RegisterDelegate, UIImagePickerControllerDelegate, UINavigationControllerDelegate>

@property (strong, nonatomic) HomeView* homeView; // 注册页面
@property (strong, nonatomic) NSString* imageUrl; // 头像上传成功后获取到的路径
@property (assign, nonatomic) NSInteger hasImage; // 是否选择了图片：0（未选择） 1（选择）
@property (strong, nonatomic) UIImage* headImage; // 图片
@property (strong, nonatomic) NSString* headImageUrl; // 远程路径

@end

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    
    _hasImage = 0; // 未选择图片
    
    //注册键盘弹出通知
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillShowNotification object:nil];
    //注册键盘隐藏通知
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(keyboardShowAndHide:) name:UIKeyboardWillHideNotification object:nil];
    
    [self createHomeView];
}

# pragma mark - 创建HomeView
- (void)createHomeView {
    self.homeView = [[HomeView alloc] initWithFrame:self.view.bounds];
    self.homeView.delegate = self;
    [self.view addSubview:self.homeView];
    
    // 渐变器
    CAGradientLayer *gradient = [CAGradientLayer layer];
    gradient.frame = self.view.bounds;
    NSMutableArray *graColor = [[NSMutableArray alloc] initWithObjects:
                                (id)[UIColor colorWithRed:250 / 250.0 green:94 / 250.0 blue:99 / 250.0 alpha:1].CGColor,
                                (id)[UIColor colorWithRed:227 / 250.0 green:40 / 250.0 blue:117 / 250.0 alpha:1].CGColor,
                                nil];
    gradient.colors = graColor;
    
    // 渐变
    gradient.startPoint = CGPointMake(0, 0);
    gradient.endPoint = CGPointMake(0, 1);
    [self.view.layer insertSublayer:gradient atIndex:0];
}

# pragma mark - 选择图片代理
-(void) selectImage {
    NSInteger type = 0; // 0 - 图库， 1 - 相机， 2 - 相册
    if ([UIImagePickerController isSourceTypeAvailable:UIImagePickerControllerSourceTypePhotoLibrary]) {
        type = UIImagePickerControllerSourceTypePhotoLibrary;
    }
    // 实例化图片来源类型
    UIImagePickerController* pick = [[UIImagePickerController alloc] init];
    pick.allowsEditing = YES; // 允许编辑
    pick.delegate = self;
    pick.sourceType = type;
    [self presentViewController:pick animated:YES completion:nil];
}

# pragma mark - 选择图片后
- (void)imagePickerController:(UIImagePickerController *)picker didFinishPickingMediaWithInfo:(NSDictionary<UIImagePickerControllerInfoKey,id> *)info {
    [picker dismissViewControllerAnimated:YES completion:nil];
    UIImage* image = [info objectForKey:UIImagePickerControllerEditedImage];
    self.homeView.headImage.image = image;
    if (image) {
        self.homeView.headLabel.hidden = YES;
        _hasImage = 1;
    }
    else {
        self.homeView.headLabel.hidden = NO;
        _hasImage = 0;
    }
    if (@available(iOS 11.0, *)) {
        self.imageUrl = [info objectForKey:UIImagePickerControllerImageURL];
    }
    else {
        self.imageUrl = [info objectForKey:UIImagePickerControllerPHAsset];
    }
    
    self.headImage = image;
}

# pragma mark - 选取图片时取消
-(void)imagePickerControllerDidCancel:(UIImagePickerController *)picker {
    [self dismissViewControllerAnimated:YES completion:nil];
}

# pragma mark - 点击空白 退出编辑
- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event {
    [self.homeView.nameTextFiled resignFirstResponder];
    [self.homeView.emailTextFiled resignFirstResponder];
    [self.homeView.passwordTextFiled resignFirstResponder];
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

# pragma mark - 按键响应
- (void)clickButton:(NSInteger)btnTag {
    switch (btnTag) {
        case 1: // 注册
            [self requestImgUrl];
            break;
        case 2: // 登录
            [self clickGotoLogin];
            break;
        default:
            break;
    }
}

# pragma mark - 上传图片
- (void)requestImgUrl {
    if (_hasImage == 1) {
        [HttpData requestImgUploadFile:self.headImage success:^(id  _Nonnull json) {
            if ([json isKindOfClass:[NSDictionary class]]) {
                if ([json[@"status"]integerValue] == 200) {
                    self.headImageUrl = json[@"data"];
                    [self clickReginster];
                }
                else {
                    [SVProgressHUD dismissWithDelay:1.5];
                    [SVProgressHUD showErrorWithStatus:json[@"message"]];
                }
            }
        } failure:^(NSError * _Nonnull err) {
            [SVProgressHUD dismissWithDelay:1.5];
            [SVProgressHUD showErrorWithStatus:@"图片上传失败!"];
        }];
    }
    else {
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showInfoWithStatus:@"请选择头像!"];
    }
}

# pragma mark - 注册
- (void)clickReginster {
        [HttpData reginsterWithUserType:@"1" withName:self.homeView.nameTextFiled.text withEmail:self.homeView.emailTextFiled.text withPhoto:self.imageUrl withPassword:self.homeView.passwordTextFiled.text seccess:^(id  _Nonnull json) {
            if ([json isKindOfClass:[NSDictionary class]]) {
                if ([json[@"status"]integerValue] == 200) {
                    [UserDefaults putUserDefaults:@"name" Value:self.homeView.nameTextFiled.text];
                    [UserDefaults putUserDefaults:@"email" Value:self.homeView.emailTextFiled.text];
                    [UserDefaults putUserDefaults:@"password" Value:self.homeView.passwordTextFiled.text];
                    
                    [UserDefaults putUserDefaults:@"isLogin" Value:@"YES"];
                    [UserDefaults putUserDefaults:@"token" Value:json[@"token"]];
                
                    NSDate* imageData = UIImageJPEGRepresentation(self.homeView.headImage.image, 1.0);
                    [UserDefaults putUserDefaults:@"image" Value:imageData];
                    
                    [SVProgressHUD dismissWithDelay:1.5];
                    [SVProgressHUD showSuccessWithStatus:@"注册成功!"];
                    
                    HomeViewController* homeVC = [[HomeViewController alloc] init];
//                    UINavigationController *navHome = [[UINavigationController alloc] initWithRootViewController:homeVC];
//                    KeyWindow.rootViewController = navHome;
                    [self.navigationController pushViewController:homeVC animated:YES];
                }
                else {
                    [SVProgressHUD dismissWithDelay:1.5];
                    [SVProgressHUD showErrorWithStatus:json[@"message"]];
                }
            }
        } failure:^(NSError * _Nonnull err) {
            [SVProgressHUD dismissWithDelay:1.5];
            [SVProgressHUD showErrorWithStatus:@"注册失败!"];
//            NSLog(@"注册失败");
        }];
}

# pragma mark - 跳转登录
- (void)clickGotoLogin {
    LoginViewController* logVC = [[LoginViewController alloc] init];
    [self.navigationController pushViewController:logVC animated:YES];
}

@end
