//
//  PersonalInfoTopView.m
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "PersonalInfoTopView.h"

@implementation PersonalInfoTopView

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.topView];
        [self.topView addSubview:self.cancelBtn];
        [self.topView addSubview:self.logoutBtn];
        [self.topView addSubview:self.saveBtn];
        [self.topView addSubview:self.titleLabel];
        self.titleLabel.hidden = YES;
    }
    return self;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    _topView.sd_layout.leftSpaceToView(self, 0).rightSpaceToView(self, 0).topSpaceToView(self, 0).bottomSpaceToView(self, 0).heightIs(60).widthIs(SCREEN_WIDTH);
    _cancelBtn.sd_layout.leftSpaceToView(_topView, 15).bottomSpaceToView(_topView, 0).heightIs(60).widthIs(60);
    _saveBtn.sd_layout.rightSpaceToView(_topView, 15).bottomSpaceToView(_topView, 0).heightIs(60).widthIs(60);
    _logoutBtn.sd_layout.rightSpaceToView(_saveBtn, 15).bottomSpaceToView(_topView, 0).heightIs(60).widthIs(60);
    _titleLabel.sd_layout.rightSpaceToView(_topView, SCREEN_WIDTH / 2 - 40).bottomSpaceToView(_topView, 0).heightIs(60).widthIs(80);
}

- (UIView*)topView {
    if (!_topView) {
        _topView = [[UIView alloc] init];
    }
    return _topView;
}

- (UIButton*)cancelBtn {
    if (!_cancelBtn) {
        _cancelBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
        _cancelBtn.tag = 1;
        [_cancelBtn setTitle:@"Cancel" forState:UIControlStateNormal];
        [_cancelBtn setTintColor:[UIColor colorWithRed:65/255.0 green:105/255.0 blue:225/255.0 alpha:1]];
        _cancelBtn.titleLabel.font = Font(19);
        [_cancelBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _cancelBtn;
}

- (UIButton*)logoutBtn {
    if (!_logoutBtn) {
        _logoutBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
        _logoutBtn.tag = 2;
        [_logoutBtn setTitle:@"Logout" forState:UIControlStateNormal];
        [_logoutBtn setTintColor:[UIColor colorWithRed:65/255.0 green:105/255.0 blue:225/255.0 alpha:1]];
        _logoutBtn.titleLabel.font = Font(19);
        [_logoutBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _logoutBtn;
}

- (UIButton*)saveBtn {
    if (!_saveBtn) {
        _saveBtn = [UIButton buttonWithType:UIButtonTypeRoundedRect];
        _saveBtn.tag = 3;
        [_saveBtn setTitle:@"Save" forState:UIControlStateNormal];
        [_saveBtn setTintColor:[UIColor colorWithRed:65/255.0 green:105/255.0 blue:225/255.0 alpha:1]];
        _saveBtn.titleLabel.font = Font(19);
        [_saveBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _saveBtn;
}

- (UILabel*)titleLabel {
    if (!_titleLabel) {
        _titleLabel = [[UILabel alloc] init];
        _titleLabel.text = @"Setting";
        _titleLabel.textColor = [UIColor blackColor];
        _titleLabel.font = BoldFont(20);
    }
    return _titleLabel;
}

- (void)clickButton:(UIButton*)btn {
    NSLog(@"%ld", btn.tag);
}

@end
