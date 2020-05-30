//
//  AddressBookTopView.m
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "AddressBookTopView.h"

@implementation AddressBookTopView

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.backBtn];
        [self addSubview:self.titleImg];
        [self addSubview:self.messageBtn];
        [self addSubview:self.listBtn];
        [self addSubview:self.lineLab];
    }
    return self;
}

- (void)layoutSubviews {
    [super layoutSubviews];
    _backBtn.sd_layout.leftSpaceToView(self, 20).topSpaceToView(self, 50).widthIs(50).heightIs(50);
    _titleImg.sd_layout.leftSpaceToView(self, SCREEN_WIDTH / 2 - 25).topSpaceToView(self, SCREEN_HEIGHT / 8 - 20).widthIs(50).heightIs(40);
    _messageBtn.sd_layout.leftSpaceToView(self, 0).bottomSpaceToView(self, 0).widthIs(SCREEN_WIDTH / 2 - 1).heightIs(80);
    _listBtn.sd_layout.rightSpaceToView(self, 0).bottomSpaceToView(self, 0).widthIs(SCREEN_WIDTH / 2 - 1).heightIs(80);
    _lineLab.sd_layout.leftSpaceToView(self, SCREEN_WIDTH / 2 - 1).bottomSpaceToView(self, 20).widthIs(2).heightIs(40);
}

- (UIButton*)backBtn {
    if (!_backBtn) {
        _backBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _backBtn.tag = 1;
        [_backBtn setImage:[UIImage imageNamed:@"friBack"] forState:UIControlStateNormal];
        [_backBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _backBtn;
}

- (UIImageView*)titleImg {
    if (!_titleImg) {
        _titleImg = [[UIImageView alloc] init];
        _titleImg.image = [UIImage imageNamed:@"friTitleimg"];
    }
    return _titleImg;
}

- (UIButton*)messageBtn {
    if (!_messageBtn) {
        _messageBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _messageBtn.tag = 2;
        [_messageBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
        [_messageBtn setTitle:@"Messages" forState:UIControlStateNormal];
        [_messageBtn setTitleColor:[UIColor grayColor] forState:UIControlStateNormal];
        [_messageBtn setTitleColor:[UIColor colorWithRed:178 / 255.0 green:34 / 255.0 blue:34 / 255.0 alpha:1] forState:UIControlStateSelected];
        _messageBtn.titleLabel.font = BoldFont(20);
        _messageBtn.selected = YES;
    }
    return _messageBtn;
}

- (UIButton*)listBtn {
    if (!_listBtn) {
        _listBtn = [UIButton buttonWithType:UIButtonTypeCustom];
        _listBtn.tag = 3;
        [_listBtn addTarget:self action:@selector(clickButton:) forControlEvents:UIControlEventTouchUpInside];
        [_listBtn setTitle:@"List" forState:UIControlStateNormal];
        [_listBtn setTitleColor:[UIColor grayColor] forState:UIControlStateNormal];
        [_listBtn setTitleColor:[UIColor colorWithRed:178 / 255.0 green:34 / 255.0 blue:34 / 255.0 alpha:1] forState:UIControlStateSelected];
        _listBtn.titleLabel.font = BoldFont(20);
    }
    return _listBtn;
}

- (UILabel*)lineLab {
    if (!_lineLab) {
        _lineLab = [[UILabel  alloc] init];
        _lineLab.backgroundColor = [UIColor colorWithRed:240 / 255.0 green:240 / 255.0 blue:240 / 255.0 alpha:1];
    }
    return _lineLab;
}


- (void)clickButton:(UIButton*)btn {
//    NSLog(@"%ld", btn.tag);
    if (btn.tag == 2) {
        btn.selected = YES;
        _listBtn.selected = NO;
    }
    if (btn.tag == 3) {
        btn.selected = YES;
        _messageBtn.selected = NO;
    }
    if ([self.delegate respondsToSelector:@selector(clickAddressBookTopViewButton:)]) {
        [self.delegate clickAddressBookTopViewButton:btn.tag];
    }
}


@end
