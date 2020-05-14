//
//  CardView.h
//  Tinder
//
//  Created by Layer on 2020/5/9.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@protocol CardViewDelegate <NSObject>

- (void) clickInfoBtn:(UIButton*)btn;

@end

@interface CardView : UIView

@property (strong, nonatomic) UIImageView* image;
@property (strong, nonatomic) UILabel* name;
@property (strong, nonatomic) UILabel* age;
@property (strong, nonatomic) UILabel* work;
@property (strong, nonatomic) UIButton* info;

@property (strong, nonatomic) UILabel* like;
@property (strong, nonatomic) UILabel* dislike;

@property (weak, nonatomic) id <CardViewDelegate> delegate;

- (void)homeArrData:(NSDictionary*) str;

@end

NS_ASSUME_NONNULL_END
