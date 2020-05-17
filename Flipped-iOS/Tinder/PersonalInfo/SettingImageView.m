//
//  SettingImageView.m
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "SettingImageView.h"

@implementation SettingImageView

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.topView];
        [self.topView addSubview:self.settingLabel];
        [self addSubview:self.image1];
        [self addSubview:self.image2];
        [self addSubview:self.image3];
        [self addSubview:self.nameLabel];
    }
    return self;
}

-(void)layoutSubviews {
    [super layoutSubviews];
    _topView.sd_layout.leftSpaceToView(self, 0).topSpaceToView(self, 0).heightIs(60).widthIs(SCREEN_WIDTH);
    _settingLabel.sd_layout.leftSpaceToView(_topView, 15).topSpaceToView(_topView, 5).rightSpaceToView(_topView, 15).heightIs(60);
    _nameLabel.sd_layout.leftSpaceToView(self, 15).rightSpaceToView(self, 15).topSpaceToView(self, SCREEN_WIDTH - 60).heightIs(60).widthIs(SCREEN_WIDTH - 30);
}

- (UIView*)topView {
    if (!_topView) {
        _topView = [[UIView alloc] init];
        _topView.backgroundColor = [UIColor whiteColor];
    }
    return _topView;
}

- (UILabel*)settingLabel {
    if (!_settingLabel) {
        _settingLabel = [[UILabel alloc] init];
        _settingLabel.text = @"Setting";
        _settingLabel.font = BoldFont(30);
    }
    return _settingLabel;
}

- (UIButton*)image1 {
    if (!_image1) {
        _image1 = [[UIButton alloc] initWithFrame:CGRectMake(10, 70, SCREEN_WIDTH / 2 - 15, SCREEN_WIDTH - 140)];
        _image1.backgroundColor = [UIColor whiteColor];
        _image1.layer.cornerRadius = 10;
        _image1.layer.masksToBounds = YES;
        _image1.userInteractionEnabled =YES;
        
        _image1.tag = 1;
        [_image1 addTarget:self action:@selector(clickImage:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _image1;
}

- (UIButton*)image2 {
    if (!_image2) {
        _image2 = [[UIButton alloc] initWithFrame:CGRectMake(SCREEN_WIDTH / 2 + 5, 70, SCREEN_WIDTH / 2 - 15, SCREEN_WIDTH / 2 - 75)];
        _image2.backgroundColor = [UIColor whiteColor];
        _image2.layer.cornerRadius = 10;
        _image2.layer.masksToBounds = YES;
        _image2.userInteractionEnabled =YES;
        
        _image2.tag = 2;
        [_image2 addTarget:self action:@selector(clickImage:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _image2;
}

- (UIButton*)image3 {
    if (!_image3) {
        _image3 = [[UIButton alloc] initWithFrame:CGRectMake(SCREEN_WIDTH / 2 + 5, SCREEN_WIDTH / 2 + 5, SCREEN_WIDTH / 2 - 15, SCREEN_WIDTH / 2 - 75)];
        _image3.backgroundColor = [UIColor whiteColor];
        _image3.layer.cornerRadius = 10;
        _image3.layer.masksToBounds = YES;
        _image3.userInteractionEnabled =YES;
        
        _image3.tag = 3;
        [_image3 addTarget:self action:@selector(clickImage:) forControlEvents:UIControlEventTouchUpInside];
        
    }
    return _image3;
}

- (UILabel*)nameLabel {
    if (!_nameLabel) {
        _nameLabel = [[UILabel alloc] initWithFrame:CGRectMake(15, SCREEN_WIDTH - 50, SCREEN_WIDTH - 15, SCREEN_WIDTH - 15)];
        _nameLabel.text = @"Name";
        _nameLabel.textColor = [UIColor blackColor];
        _nameLabel.font = BoldFont(20);
    }
    return _nameLabel;
}

- (void)clickImage:(UIButton*)btn {
    if ([_delegate respondsToSelector:@selector(clickImageButton:)]) {
        [_delegate clickImageButton:btn.tag];
    }
}

@end
