//
//  InfoBottonView.m
//  Tinder
//
//  Created by Layer on 2020/5/14.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "InfoBottonView.h"

@implementation InfoBottonView

#pragma mark -- 初始化
- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.nameLabel];
//        [self addSubview:self.ageLabel];
        [self addSubview:self.jobLabel];
        [self addSubview:self.infoLabel];
        [self addSubview:self.disLikeBtn];
        [self addSubview:self.collectBtn];
        [self addSubview:self.likeBtn];
    }
    return self;
}

#pragma mark -- 位置
- (void)layoutSubviews {
    [super layoutSubviews];
    _nameLabel.sd_layout.leftSpaceToView(self, 20).topSpaceToView(self, 20).rightSpaceToView(self, 50).heightIs(40);
//    _ageLabel.sd_layout.leftSpaceToView(_nameLabel, 0).topSpaceToView(self, 20).rightSpaceToView(self, 20).heightIs(40);
    _jobLabel.sd_layout.leftSpaceToView(self, 20).topSpaceToView(_nameLabel, 5).rightSpaceToView(self, 110).heightIs(40);
    _infoLabel.sd_layout.leftSpaceToView(self, 20).topSpaceToView(_jobLabel, 20).rightSpaceToView(self, 110).heightIs(40);
}

#pragma mark -- 名字
- (UILabel*)nameLabel {
    if (!_nameLabel) {
        _nameLabel = [[UILabel alloc] init];
        _nameLabel.text = @"Crayon Shinchan";
        _nameLabel.font = BoldFont(35);
    }
    return _nameLabel;
}

#pragma mark -- 年龄
- (UILabel*)ageLabel {
    if (!_ageLabel) {
        _ageLabel = [[UILabel alloc] init];
        _ageLabel.text = @"5";
        _ageLabel.font = BoldFont(25);
    }
    return _ageLabel;
}

#pragma mark -- 工作
- (UILabel*)jobLabel {
    if (!_jobLabel) {
        _jobLabel = [[UILabel alloc] init];
        _jobLabel.text = @"Student";
        _jobLabel.font = BoldFont(25);
    }
    return _jobLabel;
}

#pragma mark -- 信息
- (UILabel*)infoLabel {
    if (!_infoLabel) {
        _infoLabel = [[UILabel alloc] init];
        _infoLabel.text = @"Nice to meet you!";
        _infoLabel.font = Font(20);
    }
    return _infoLabel;
}

#pragma mark -- 按键
- (UIButton*)disLikeBtn {
    if (!_disLikeBtn) {
        _disLikeBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _disLikeBtn.frame = CGRectMake(SCREEN_WIDTH/6, SCREEN_HEIGHT/2/3*2 - 50, 80, 80);
        [_disLikeBtn setBackgroundColor: [UIColor colorWithRed:245/255.0 green:245/255.0 blue:245/255.0 alpha:1]];
        _disLikeBtn.tag = 1;
        _disLikeBtn.layer.backgroundColor = [UIColor colorWithRed:240/255.0 green:240/255.0 blue:240/255.0 alpha:1].CGColor;
        _disLikeBtn.layer.cornerRadius = 40;
        _disLikeBtn.layer.masksToBounds = YES;
        [_disLikeBtn setBackgroundImage:[UIImage imageNamed:@"inforNoLike"] forState:UIControlStateNormal];
        // setBackgroundImg,会适应大小，setImage不会适应大小
        [_disLikeBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _disLikeBtn;
}


- (UIButton*)collectBtn {
    if (!_collectBtn) {
        _collectBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _collectBtn.frame = CGRectMake(SCREEN_WIDTH/2 - 25, SCREEN_HEIGHT/2/3*2 - 35, 50, 50);
        [_collectBtn setBackgroundColor: [UIColor colorWithRed:245/255.0 green:245/255.0 blue:245/255.0 alpha:1]];
        _collectBtn.tag = 2;
        _collectBtn.layer.backgroundColor = [UIColor colorWithRed:240/255.0 green:240/255.0 blue:240/255.0 alpha:1].CGColor;
        _collectBtn.layer.cornerRadius = 25;
        _collectBtn.layer.masksToBounds = YES;
        [_collectBtn setBackgroundImage:[UIImage imageNamed:@"inforSC"] forState:UIControlStateNormal];
        [_collectBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _collectBtn;
}

- (UIButton*)likeBtn {
    if (!_likeBtn) {
        _likeBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _likeBtn.frame = CGRectMake(SCREEN_WIDTH - SCREEN_WIDTH/6 - 80, SCREEN_HEIGHT/2/3*2 - 50, 80, 80);
        [_likeBtn setBackgroundColor: [UIColor colorWithRed:245/255.0 green:245/255.0 blue:245/255.0 alpha:1]];
        _likeBtn.tag = 3;
        _likeBtn.layer.backgroundColor = [UIColor colorWithRed:240/255.0 green:240/255.0 blue:240/255.0 alpha:1].CGColor;
        _likeBtn.layer.cornerRadius = 40;
        _likeBtn.layer.masksToBounds = YES;
        [_likeBtn setBackgroundImage:[UIImage imageNamed:@"inforLike"] forState:UIControlStateNormal];
        [_likeBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _likeBtn;
}

- (void)clickButton:(UIButton*)btn {
    if ([_delegate respondsToSelector:@selector(clickinfoButton:)]) {
        [_delegate clickinfoButton:btn.tag];
    }
}

@end
