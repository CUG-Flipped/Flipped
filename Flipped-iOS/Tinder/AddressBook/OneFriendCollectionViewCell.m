//
//  OneFriendCollectionViewCell.m
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "OneFriendCollectionViewCell.h"

@implementation OneFriendCollectionViewCell

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        _headImg = [[UIImageView alloc] init];
        _headImg.backgroundColor = [UIColor colorWithRed:245 / 255.0 green:245 / 255.0 blue:245 / 255.0 alpha:1];
        _headImg.layer.cornerRadius = (SCREEN_WIDTH / 4 - 10) / 2;
        _headImg.layer.masksToBounds = YES;
        _headImg.image = [UIImage imageNamed:@"navLeftImg"];
        [self.contentView addSubview: _headImg];
        _headImg.sd_layout.leftSpaceToView(self.contentView, 5).topSpaceToView(self.contentView, 5).rightSpaceToView(self.contentView, 5).heightIs(SCREEN_WIDTH / 4 - 10);
        
        _nameLab = [[UILabel alloc] init];
        _nameLab.numberOfLines = 0;
        _nameLab.text = @"Crayon Shinchan";
        _nameLab.textAlignment = NSTextAlignmentCenter;
        _nameLab.font = BoldFont(13);
        _nameLab.textColor = [UIColor colorWithRed:30 / 255.0 green:30 / 255.0 blue:30 / 255.0 alpha:1];
        [self.contentView addSubview:_nameLab];
        _nameLab.sd_layout.leftSpaceToView(self.contentView, 5).rightSpaceToView(self.contentView, 5).topSpaceToView(_headImg, 5).bottomSpaceToView(self.contentView, 0);
    }
    return self;
}

@end
