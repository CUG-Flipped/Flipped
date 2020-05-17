//
//  SettingImageView.h
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol SettingImageViewDelegate <NSObject>

- (void)clickImageButton:(NSInteger)tag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface SettingImageView : UIView

@property(strong, nonatomic) UIView *topView;
@property(strong, nonatomic) UILabel *settingLabel;

@property(strong, nonatomic) UIButton *image1;
@property(strong, nonatomic) UIButton *image2;
@property(strong, nonatomic) UIButton *image3;

@property(strong, nonatomic) UILabel *nameLabel;

@property(weak, nonatomic) id<SettingImageViewDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
