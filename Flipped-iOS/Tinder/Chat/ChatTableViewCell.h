//
//  ChatTableViewCell.h
//  Tinder
//
//  Created by Layer on 2020/5/30.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@interface ChatTableViewCell : UITableViewCell

@property (strong, nonatomic)UIButton *backButton;
@property (strong, nonatomic)UILabel *messageLabel;

@end

NS_ASSUME_NONNULL_END
