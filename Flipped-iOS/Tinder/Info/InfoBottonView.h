//
//  InfoBottonView.h
//  Tinder
//
//  Created by Layer on 2020/5/14.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@interface InfoBottonView : UIView

@property (strong, nonatomic) UILabel *nameLabel;
@property (strong, nonatomic) UILabel *ageLabel;
@property (strong, nonatomic) UILabel *jobLabel;
@property (strong, nonatomic) UILabel *infoLabel;

@property (strong, nonatomic) UIButton *disLikeBtn;
@property (strong, nonatomic) UIButton *likeBtn;
@property (strong, nonatomic) UIButton *collectBtn;;
@end

NS_ASSUME_NONNULL_END
