//
//  InfoViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/10.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "InfoViewController.h"

@interface InfoViewController () <infoButtonViewDelegate>

@property (strong, nonatomic) InfoImageView *infoImage;
@property (strong, nonatomic) InfoBottonView *infoButton;

@property (strong, nonatomic) NSMutableArray *dataArray;
@property (strong, nonatomic) NSDictionary *listDic;

@end

@implementation InfoViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    self.view.backgroundColor = [UIColor whiteColor];
    self.dataArray = [[NSMutableArray alloc] init];
    [self requesrData];
//    [self creatInfoImage];
    [self creatInfoButton];
    [self creakBackButton];
}

- (void)creatInfoImage {
//    NSMutableArray *arr = [[NSMutableArray alloc] initWithObjects: @"background_1", @"background_2", @"background_3", @"background_4",nil];
//    _infoImage = [[InfoImageView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 2) addImageArr:arr];
    _infoImage = [[InfoImageView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 2) addImageArr:self.dataArray];
    [self.view addSubview: _infoImage];
}

- (void)creatInfoButton {
    _infoButton = [[InfoBottonView alloc] initWithFrame:CGRectMake(0, SCREEN_HEIGHT / 2, SCREEN_WIDTH, SCREEN_HEIGHT/ 2)];
    _infoButton.backgroundColor = [UIColor whiteColor];
    _infoButton.delegate = self;
    [self.view addSubview: _infoButton];
}

- (void)creakBackButton {
    UIButton *btn =[UIButton buttonWithType:UIButtonTypeCustom];
//    btn.frame = CGRectMake(SCREEN_WIDTH - 90, SCREEN_HEIGHT / 2 - 25, 50, 50);
    btn.frame = CGRectMake(SCREEN_WIDTH - 90, SCREEN_HEIGHT / 2 + 25, 50, 50);
    [btn setBackgroundImage:[UIImage imageNamed:@"infoBackDown"] forState:UIControlStateNormal];
    [btn addTarget:self action:@selector(clickBackButton) forControlEvents:UIControlEventTouchUpInside];
    [self.view addSubview:btn];
}

- (void)clickBackButton {
    [self dismissViewControllerAnimated:YES completion:nil];
}

- (void)clickinfoButton:(NSInteger)tag {
    NSLog(@"%ld", tag);
}

- (void)requesrData {
    kWeakSelf(self);
    [HttpData requestHomeInfoToken:[UserDefaults getUserDefaults:@"token"] user_id:@"36" success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"] integerValue] == 200) {
                weakself.listDic = json[@"data"];
                for (int i = 0; i < [weakself.listDic[@"photo"] count]; i++) {
                    [weakself.dataArray addObject:weakself.listDic[@"photo"][i]];
                }
                [self creatInfoImage];
                weakself.infoButton.nameLabel.text = [NSString stringWithFormat:@"%@ %@", weakself.listDic[@"user_name"],weakself.listDic[@"age"]];
                weakself.infoButton.jobLabel.text = weakself.listDic[@"work"];
                if (!(weakself.listDic[@"bio"] == nil ||[weakself.listDic[@"bio"] isEqual:[NSNull null]])) {
                    weakself.infoButton.infoLabel.text = weakself.listDic[@"bio"];
                }
            }
            else {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:json[@"message"]];
            }
        }
    } failure:^(NSError * _Nonnull err) {
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showErrorWithStatus:@"Error"];
    }];
}

@end
