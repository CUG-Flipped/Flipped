//
//  HomeTopView.m
//  Tinder
//
//  Created by Layer on 2020/5/8.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "HomeTopView.h"

@implementation HomeTopView

# pragma mark - 初始化
- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview: self.leftBtn];
        [self addSubview: self.rightBtn];
        [self addSubview: self.topImage];
    }
    return self;
}

# pragma mark - 位置
- (void)layoutSubviews {
    [super layoutSubviews];
    _leftBtn.sd_layout.leftSpaceToView(self, 20).bottomSpaceToView(self, 15).widthIs(40).heightIs(40);
    _rightBtn.sd_layout.rightSpaceToView(self, 20).bottomSpaceToView(self, 15).widthIs(40).heightIs(40);
    _topImage.sd_layout.leftSpaceToView(self, SCREEN_WIDTH / 2 - 15).bottomSpaceToView(self, 15).widthIs(30).heightIs(40);
}

# pragma mark - 左按键
- (UIButton*)leftBtn {
    if (!_leftBtn) {
        _leftBtn = [UIButton buttonWithType:UIButtonTypeCustom]; // 自定义类型
        [_leftBtn setImage:[UIImage imageNamed:@"navLeftImg"] forState:UIControlStateNormal]; // 背景
        
        _leftBtn.tag = 1;
        [_leftBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _leftBtn;
}

# pragma mark - 右按键
- (UIButton*)rightBtn {
    if (!_rightBtn) {
        _rightBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        [_rightBtn setImage:[UIImage imageNamed:@"navRightImg"] forState:UIControlStateNormal];
        _rightBtn.tag = 2;
        [_rightBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _rightBtn;
}

# pragma mark - 图标
- (UIImageView*)topImage {
    if (!_topImage) {
        _topImage = [[UIImageView alloc] init];
        _topImage.image = [UIImage imageNamed:@"navTitleImg"];
    }
    return _topImage;
}

# pragma mark - 按键响应
- (void)clickButton:(UIButton*)btn {
    if ([_delegate respondsToSelector:@selector(clickButton:)]) {
        [_delegate clickButton:btn.tag];
    }
}

@end
