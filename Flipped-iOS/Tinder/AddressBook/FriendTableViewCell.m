//
//  FriendTableViewCell.m
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "FriendTableViewCell.h"

@implementation FriendTableViewCell

- (void)awakeFromNib {
    [super awakeFromNib];
    // Initialization code
}

- (instancetype)initWithStyle:(UITableViewCellStyle)style reuseIdentifier:(NSString *)reuseIdentifier {
    if (self = [super initWithStyle:style reuseIdentifier:reuseIdentifier]) {
        _headImg = [[UIImageView alloc] init];
        _headImg.backgroundColor = [UIColor colorWithRed:245 / 255.0 green:245 / 255.0 blue:245 / 255.0 alpha:1];
        _headImg.layer.cornerRadius = 40;
        _headImg.layer.masksToBounds = YES;
        [self.contentView addSubview:_headImg];
        _headImg.sd_layout.leftSpaceToView(self.contentView, 10).topSpaceToView(self.contentView, 20).widthIs(80).heightIs(80);
        
        _nameLab = [[UILabel alloc] init];
        _nameLab.font = BoldFont(24);
        [self.contentView addSubview:_nameLab];
        _nameLab.sd_layout.leftSpaceToView(_headImg, 15).topSpaceToView(self.contentView, 30).rightSpaceToView(self.contentView, 15).heightIs(35);
        
        _messageLab = [[UILabel alloc] init];
        _messageLab.font = BoldFont(18);
        _messageLab.numberOfLines = 0; // 自动换行
        _messageLab.textColor = [UIColor colorWithRed:180 / 255.0 green:180 / 255.0 blue:180 / 255.0 alpha:1];
        [self.contentView addSubview:_messageLab];
        _messageLab.sd_layout.leftSpaceToView(_headImg, 15).bottomSpaceToView(self.contentView, 30).rightSpaceToView(self.contentView, 15).heightIs(35);
        
    }
    return self;
}


- (void)setSelected:(BOOL)selected animated:(BOOL)animated {
    [super setSelected:selected animated:animated];

    // Configure the view for the selected state
}


@end
