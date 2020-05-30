//
//  ChatTopView.m
//  Tinder
//
//  Created by Layer on 2020/5/24.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "ChatTopView.h"

@implementation ChatTopView

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.backBtn];
        [self addSubview:self.headImg];
        [self addSubview:self.nameLabel];
        [self addSubview:self.flagBtn];
    }
    return self;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    _backBtn.sd_layout.leftSpaceToView(self, 10).topSpaceToView(self, SCREEN_HEIGHT / 14).widthIs(30).heightIs(30);
    _headImg.sd_layout.leftSpaceToView(self, SCREEN_WIDTH / 2 - 20).topSpaceToView(self, 50).widthIs(40).heightIs(40);
    _nameLabel.sd_layout.leftSpaceToView(self, 10).rightSpaceToView(self, 10).topSpaceToView(_headImg, 10).bottomSpaceToView(self, 0);
    _flagBtn.sd_layout.rightSpaceToView(self, 10).topSpaceToView(self, SCREEN_HEIGHT / 14).widthIs(30).heightIs(30);
}

- (UIButton*)backBtn {
    if (!_backBtn) {
        _backBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _backBtn.tag = 1;
//        [_backBtn setImage:[UIImage imageNamed:@"sepBack"] forState:UIControlStateNormal];
        [_backBtn setBackgroundImage:[UIImage imageNamed:@"sepBack"] forState:UIControlStateNormal];
        [_backBtn addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _backBtn;
}

- (UIImageView*)headImg {
    if (!_headImg) {
        _headImg = [[UIImageView alloc] init];
        _headImg.image = [UIImage imageNamed:@"navLeftImg"];
        _headImg.layer.cornerRadius = 20;
        _headImg.layer.masksToBounds = YES;
    }
    return _headImg;
}

- (UILabel*)nameLabel {
    if (!_nameLabel) {
        _nameLabel = [[UILabel alloc] init];
        _nameLabel.text = @"Crayon Shinchan";
        _nameLabel.font = BoldFont(15);
        _nameLabel.textAlignment = NSTextAlignmentCenter;
    }
    return _nameLabel;
}

- (UIButton*)flagBtn {
    if (!_flagBtn) {
        _flagBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _flagBtn.tag = 2;
//        [_flagBtn setImage:[UIImage imageNamed:@"sepRightImg"] forState:UIControlStateNormal];
        [_flagBtn setBackgroundImage:[UIImage imageNamed:@"sepRightImg"] forState:UIControlStateNormal];
        [_flagBtn addTarget:self action:@selector(clickBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _flagBtn;
}

- (void)clickBtn:(UIButton*)btn {
//    NSLog(@"%ld", btn.tag);
    if ([self.delegate respondsToSelector:@selector(chatTopViewButtonClick:)]) {
        [self.delegate chatTopViewButtonClick:btn.tag];
    }
}


@end
