//
//  ChatTableViewCell.m
//  Tinder
//
//  Created by Layer on 2020/5/30.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "ChatTableViewCell.h"

@implementation ChatTableViewCell

- (void)awakeFromNib {
    [super awakeFromNib];
    // Initialization code
}

- (instancetype)initWithStyle:(UITableViewCellStyle)style reuseIdentifier:(NSString *)reuseIdentifier {
    if (self = [super initWithStyle:style reuseIdentifier:reuseIdentifier]) {
        
        _backButton = [UIButton buttonWithType:UIButtonTypeCustom];
        _backButton.frame = CGRectZero;
        _backButton.layer.cornerRadius = 10;
        _backButton.layer.masksToBounds = YES;
        _backButton.userInteractionEnabled = NO;
        [self.contentView addSubview:_backButton];
        
        _messageLabel = [[UILabel alloc] init];
        _messageLabel.font = Font(16);
        _messageLabel.numberOfLines = 0;
        [_backButton addSubview:_messageLabel];
        _messageLabel.sd_layout.leftSpaceToView(_backButton, 10).topSpaceToView(_backButton, 5).rightSpaceToView(_backButton, 10).bottomSpaceToView(_backButton, 5);
        
    }
    return self;
}

- (void)setSelected:(BOOL)selected animated:(BOOL)animated {
    [super setSelected:selected animated:animated];

    // Configure the view for the selected state
}

@end
