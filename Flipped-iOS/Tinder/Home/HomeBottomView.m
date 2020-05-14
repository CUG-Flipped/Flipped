//
//  HomeBottomView.m
//  Tinder
//
//  Created by Layer on 2020/5/8.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "HomeBottomView.h"

@implementation HomeBottomView

# pragma mark - 初始化
- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview: self.btn_1];
        [self addSubview: self.btn_2];
        [self addSubview: self.btn_3];
        [self addSubview: self.btn_4];
        [self addSubview: self.btn_5];
    }
    return self;
}

# pragma mark - 位置
- (void)layoutSubviews {
    [super layoutSubviews];
    _btn_1.sd_layout.leftSpaceToView(self, 20).topSpaceToView(self, 20).widthIs(SCREEN_WIDTH / 5 - 40).heightIs(SCREEN_WIDTH / 5 - 40);
    _btn_2.sd_layout.leftSpaceToView(_btn_1, 30).topSpaceToView(self, 10).widthIs(SCREEN_WIDTH / 5 - 20).heightIs(SCREEN_WIDTH / 5 - 20);
    _btn_3.sd_layout.leftSpaceToView(self, SCREEN_WIDTH / 2- (SCREEN_WIDTH / 5 - 40) / 2).topSpaceToView(self, 20).widthIs(SCREEN_WIDTH / 5 - 40).heightIs(SCREEN_WIDTH / 5 - 40);
    _btn_5.sd_layout.rightSpaceToView(self, 20).topSpaceToView(self, 20).widthIs(SCREEN_WIDTH / 5 - 40).heightIs(SCREEN_WIDTH / 5 - 40);
    _btn_4.sd_layout.rightSpaceToView(_btn_5, 30).topSpaceToView(self, 10).widthIs(SCREEN_WIDTH / 5 - 20).heightIs(SCREEN_WIDTH / 5 - 20);
    
}

# pragma mark - 按键
- (UIButton*)btn_1 {
    if (!_btn_1) {
        _btn_1 = [UIButton buttonWithType:UIButtonTypeCustom];
        _btn_1.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _btn_1.layer.cornerRadius = (SCREEN_WIDTH / 5 - 40) / 2.0;
        _btn_1.layer.masksToBounds = YES;
        [_btn_1 setImage:[UIImage imageNamed:@"firstImg"] forState:UIControlStateNormal];
        
        _btn_1.tag = 1;
        [_btn_1 addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _btn_1;
}

- (UIButton*)btn_2 {
    if (!_btn_2) {
        _btn_2 = [UIButton buttonWithType:UIButtonTypeCustom];
        _btn_2.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _btn_2.layer.cornerRadius = (SCREEN_WIDTH / 5 - 20) / 2.0;
        _btn_2.layer.masksToBounds = YES;
        [_btn_2 setImage:[UIImage imageNamed:@"twoImg"] forState:UIControlStateNormal];
        
        _btn_2.tag = 2;
        [_btn_2 addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _btn_2;
}

- (UIButton*)btn_3 {
    if (!_btn_3) {
        _btn_3 = [UIButton buttonWithType:UIButtonTypeCustom];
        _btn_3.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _btn_3.layer.cornerRadius = (SCREEN_WIDTH / 5 - 40) / 2.0;
        _btn_3.layer.masksToBounds = YES;
        [_btn_3 setImage:[UIImage imageNamed:@"threeImg"] forState:UIControlStateNormal];
        
        _btn_3.tag = 3;
        [_btn_3 addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _btn_3;
}

- (UIButton*)btn_4 {
    if (!_btn_4) {
        _btn_4 = [UIButton buttonWithType:UIButtonTypeCustom];
        _btn_4.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _btn_4.layer.cornerRadius = (SCREEN_WIDTH / 5 - 20) / 2.0;
        _btn_4.layer.masksToBounds = YES;
        [_btn_4 setImage:[UIImage imageNamed:@"fourImg"] forState:UIControlStateNormal];
        
        _btn_4.tag = 4;
        [_btn_4 addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _btn_4;
}

- (UIButton*)btn_5 {
    if (!_btn_5) {
        _btn_5 = [UIButton buttonWithType:UIButtonTypeCustom];
        _btn_5.backgroundColor = [UIColor colorWithRed:245 / 250.0 green:245 / 250.0 blue:245 / 250.0 alpha:1];
        _btn_5.layer.cornerRadius = (SCREEN_WIDTH / 5 - 40) / 2.0;
        _btn_5.layer.masksToBounds = YES;
        [_btn_5 setImage:[UIImage imageNamed:@"fiveImg"] forState:UIControlStateNormal];
        
        _btn_5.tag = 5;
        [_btn_5 addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _btn_5;
}

# pragma mark - 按键响应
- (void)clickButton:(UIButton*)btn {
    if ([_delegate respondsToSelector:@selector(homeBottomClickBtn:)]) {
        [_delegate homeBottomClickBtn:btn];
    }
}


@end
