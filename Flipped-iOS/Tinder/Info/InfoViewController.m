//
//  InfoViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/10.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "InfoViewController.h"

@interface InfoViewController ()

@property (strong, nonatomic) InfoImageView *infoImage;
@property (strong, nonatomic) InfoBottonView *infoButton;

@end

@implementation InfoViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    self.view.backgroundColor = [UIColor whiteColor];
    [self creatInfoImage];
    [self creatInfoButton];
}

- (void)creatInfoImage {
    NSMutableArray *arr = [[NSMutableArray alloc] initWithObjects: @"background_1", @"background_2", @"background_3", @"background_4",nil];
    _infoImage = [[InfoImageView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT / 2) addImageArr:arr];
    [self.view addSubview: _infoImage];
}

- (void)creatInfoButton {
    _infoButton = [[InfoBottonView alloc] initWithFrame:CGRectMake(0, SCREEN_HEIGHT / 2, SCREEN_WIDTH, SCREEN_HEIGHT/ 2)];
    _infoButton.backgroundColor = [UIColor whiteColor];
    [self.view addSubview: _infoButton];
}
/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender {
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
