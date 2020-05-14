//
//  HomeTopView.h
//  Tinder
//
//  Created by Layer on 2020/5/8.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol HomeTopDelegate <NSObject>

- (void)clickButton:(NSInteger)btnTag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface HomeTopView : UIView

@property (strong, nonatomic) UIButton* leftBtn;
@property (strong, nonatomic) UIButton* rightBtn;
@property (strong, nonatomic) UIImageView* topImage;

@property (weak, nonatomic) id<HomeTopDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
