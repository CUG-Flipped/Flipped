//
//  PersonalInfoTopView.h
//  Tinder
//
//  Created by Layer on 2020/5/16.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol PersonalInfoTopViewDelegate <NSObject>

- (void)clickPersonalInfoTopViewButton:(NSInteger)tag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface PersonalInfoTopView : UIView

@property (strong, nonatomic) UIView *topView;
@property (strong, nonatomic) UILabel *titleLabel;
@property (strong, nonatomic) UIButton *cancelBtn;
@property (strong, nonatomic) UIButton *logoutBtn;
@property (strong, nonatomic) UIButton *saveBtn;

@property (weak, nonatomic) id<PersonalInfoTopViewDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
