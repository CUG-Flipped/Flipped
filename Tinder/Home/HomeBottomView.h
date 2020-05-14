//
//  HomeBottomView.h
//  Tinder
//
//  Created by Layer on 2020/5/8.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@protocol HomeBottomDelegate <NSObject>

- (void)homeBottomClickBtn:(UIButton*)btn;

@end

@interface HomeBottomView : UIView

@property (strong, nonatomic) UIButton* btn_1;
@property (strong, nonatomic) UIButton* btn_2;
@property (strong, nonatomic) UIButton* btn_3;
@property (strong, nonatomic) UIButton* btn_4;
@property (strong, nonatomic) UIButton* btn_5;

@property (weak, nonatomic) id<HomeBottomDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
