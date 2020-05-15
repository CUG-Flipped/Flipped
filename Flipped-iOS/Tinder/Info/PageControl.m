//
//  PageControl.m
//  Tinder
//
//  Created by Layer on 2020/5/14.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "PageControl.h"

#define dotW 8
#define activeDotW 100
#define margin 6

@implementation PageControl

- (void)layoutSubviews{
    [super layoutSubviews];
    //遍历subview,设置圆点frame
    for (int i = 0; i< [self.subviews count]; i++){
        UIImageView* dot = [self.subviews objectAtIndex:i];
        [dot setFrame:CGRectMake(self.bounds.size.width/self.subviews.count * i + 2 , dot.frame.origin.y, self.bounds.size.width/self.subviews.count-2, 4)];
    }
}

- (void) setCurrentPage:(NSInteger)page{
    [super setCurrentPage:page];
    for (NSUInteger subviewIndex = 0; subviewIndex < [self.subviews count]; subviewIndex++){
        UIImageView* subview = [self.subviews objectAtIndex:subviewIndex];
        CGSize size;
        size.height = 3;
        size.width = self.bounds.size.width/self.subviews.count;
        [subview setFrame:CGRectMake(subview.frame.origin.x, subview.frame.origin.y, size.width,size.height)];
    }
}

@end
