//
//  CardView.m
//  Tinder
//
//  Created by Layer on 2020/5/9.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "CardView.h"

@implementation CardView

- (instancetype)initWithFrame:(CGRect)frame
{
    self = [super initWithFrame:frame];
    if (self) {
        [self addSubview:self.image];
        [self addSubview:self.name];
        [self addSubview:self.age];
        [self addSubview:self.work];
        [self addSubview:self.info];
        [self addSubview:self.like];
        [self addSubview:self.dislike];
    }
    return self;
}

- (UIImageView*)image {
    if (!_image) {
        _image = [[UIImageView alloc] initWithFrame:self.bounds];
        _image.image = [UIImage imageNamed: @"background_1"];
    }
    return _image;
}

- (UILabel*)name {
    if (!_name) {
        _name = [[UILabel alloc] initWithFrame:CGRectMake(15, self.frame.size.height - 110, 250, 50)];
        _name.text = @"Crayon Shinchan";
        _name.textColor = [UIColor whiteColor];
        _name.font = BoldFont(30);
    }
    return _name;
}

- (UILabel*)age {
    if (!_age) {
        _age = [[UILabel alloc] initWithFrame:CGRectMake(275, self.frame.size.height - 106, 100, 50)];
        _age.text = @"5";
        _age.textColor = [UIColor whiteColor];
        _age.font = BoldFont(20);
    }
    return _age;
}

- (UILabel*)work {
    if (!_work) {
        _work = [[UILabel alloc] initWithFrame:CGRectMake(15, self.frame.size.height - 80, 300, 50)];
        _work.text = @"Student";
        _work.textColor = [UIColor whiteColor];
        _work.font = BoldFont(20);
    }
    return _work;
}

- (UIButton*)info {
    if (!_info) {
        _info = [UIButton buttonWithType:UIButtonTypeCustom];
        _info.frame = CGRectMake(self.frame.size.width - 50, self.frame.size.height - 55, 36, 36);
        [_info setTitle:@"i" forState:UIControlStateNormal];
        _info.titleLabel.font = BoldFont(28);
        [_info setTitleColor:[UIColor blackColor] forState:UIControlStateNormal];
        [_info setBackgroundColor:[UIColor whiteColor]];
        _info.layer.cornerRadius = 18;
        _info.layer.masksToBounds =YES;
        [_info addTarget:self action:@selector(clickInfoBtn:) forControlEvents:UIControlEventTouchUpInside];
    }
    return _info;
}

- (UILabel*)like {
    if (!_like) {
        _like = [[UILabel alloc] initWithFrame:CGRectMake(15, 50, 125, 60)];
        _like.layer.cornerRadius = 5;
        _like.layer.masksToBounds = YES;
        _like.layer.borderWidth = 5;
        _like.layer.borderColor = [UIColor colorWithRed:18 / 255.0 green:238 / 255.0 blue:148 / 255.0 alpha:1].CGColor;
        _like.text = @"LIKE";
        _like.textAlignment = NSTextAlignmentCenter; // 居中
        _like.font = BoldFont(40);
        _like.textColor = [UIColor colorWithRed:18 / 255.0 green:238 / 255.0 blue:148 / 255.0 alpha:1];
        _like.transform = CGAffineTransformMakeRotation(100);
        _like.alpha = 0; // 隐藏
    }
    return _like;
}

- (UILabel*)dislike {
    if (!_dislike) {
        _dislike = [[UILabel alloc] initWithFrame:CGRectMake(self.frame.size.width - 140, 50, 125, 60)];
        _dislike.layer.cornerRadius = 5;
        _dislike.layer.masksToBounds = YES;
        _dislike.layer.borderWidth = 5;
        _dislike.layer.borderColor = [UIColor colorWithRed:205 / 255.0 green:92/ 255.0 blue:92 / 255.0 alpha:1].CGColor;
        _dislike.text = @"NOPE";
        _dislike.textAlignment = NSTextAlignmentCenter; // 居中
        _dislike.font = BoldFont(40);
        _dislike.textColor = [UIColor colorWithRed:205 / 255.0 green:92/ 255.0 blue:92 / 255.0 alpha:1];
        _dislike.transform = CGAffineTransformMakeRotation(-100);
        _dislike.alpha = 0; // 隐藏
    }
    return _dislike;
}

- (void)homeArrData:(NSDictionary*)dic {
    [self.image sd_setImageWithURL:[NSURL URLWithString:dic[@"photo"]] placeholderImage:nil];
    self.name.text = dic[@"user_name"];
    self.age.text = dic[@"user_age"];
    self.work.text = dic[@"user_work"];
    
    CGFloat nameWidth = [[NSString alloc]initWithString:dic[@"user_name"]].length * 18;
//    NSLog(@"%f",nameWidth);
    _name.sd_layout.widthIs(nameWidth);
    _age.sd_layout.leftSpaceToView(_name, 0);
}

- (void)clickInfoBtn:(UIButton*)btn {
    if ([_delegate respondsToSelector:@selector(clickInfoBtn:)]) {
        [_delegate clickInfoBtn:btn];
    }
}

@end
