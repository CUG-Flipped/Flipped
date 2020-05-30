//
//  ChatTopView.h
//  Tinder
//
//  Created by Layer on 2020/5/24.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol ChatTopViewDelegate <NSObject>

- (void)chatTopViewButtonClick:(NSInteger)tag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface ChatTopView : UIView

@property (strong, nonatomic) UIButton *backBtn;
@property (strong, nonatomic) UIImageView *headImg;
@property (strong, nonatomic) UILabel *nameLabel;
@property (strong, nonatomic) UIButton *flagBtn;

@property (weak, nonatomic) id<ChatTopViewDelegate> delegate;

@end

NS_ASSUME_NONNULL_END

