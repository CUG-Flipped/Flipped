//
//  SetSlideCellTableViewCell.h
//  Tinder
//
//  Created by Layer on 2020/5/17.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol SetBottomSlideDelegate <NSObject>

- (void)clickSlideValueTag:(NSInteger)tag withValue:(NSString*)value;   // slide值

@end

NS_ASSUME_NONNULL_BEGIN

@interface SetSlideCellTableViewCell : UITableViewCell

@property (nonatomic,strong)UILabel *minLabel;        // 最小值
@property (nonatomic,strong)UILabel *maxLabel;        // 最大值

@property (nonatomic,strong)UISlider *minSlide;        // 最小值slide
@property (nonatomic,strong)UISlider *maxslide;        // 最大值slide

@property (nonatomic,weak)id <SetBottomSlideDelegate> delegate;        // 代理

-(void)requestMinVale:(NSString *)minValue maxValue:(NSString *)maxValue;

@end

NS_ASSUME_NONNULL_END
