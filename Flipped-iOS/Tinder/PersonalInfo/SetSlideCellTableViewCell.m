//
//  SetSlideCellTableViewCell.m
//  Tinder
//
//  Created by Layer on 2020/5/17.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "SetSlideCellTableViewCell.h"

@implementation SetSlideCellTableViewCell

- (void)awakeFromNib {
    [super awakeFromNib];
    // Initialization code
}

- (void)setSelected:(BOOL)selected animated:(BOOL)animated {
    [super setSelected:selected animated:animated];

    // Configure the view for the selected state
}


- (instancetype)initWithStyle:(UITableViewCellStyle)style reuseIdentifier:(NSString *)reuseIdentifier {
    if (self = [super initWithStyle:style reuseIdentifier:reuseIdentifier]) {
        _minLabel = [[UILabel alloc]init];
        _minLabel.backgroundColor = [UIColor whiteColor];
        _minLabel.textAlignment = NSTextAlignmentCenter;
        _minLabel.font = [UIFont fontWithName:@"Arial-BoldMT" size:16];
        [self.contentView addSubview:_minLabel];
        _minLabel.sd_layout.leftSpaceToView(self.contentView, 10).topSpaceToView(self.contentView, 10).widthIs(100).heightIs(55);
        
        _maxLabel = [[UILabel alloc]init];
        _maxLabel.backgroundColor = [UIColor whiteColor];
        _maxLabel.textAlignment = NSTextAlignmentCenter;
        _maxLabel.font = [UIFont fontWithName:@"Arial-BoldMT" size:16];
        [self.contentView addSubview:_maxLabel];
        _maxLabel.sd_layout.leftSpaceToView(self.contentView, 10).topSpaceToView(_minLabel, 10).widthIs(100).heightIs(55);
        
        _minSlide = [[UISlider alloc]init];
        _minSlide.minimumValue = 16.0;  // 滑块可以滑动到最小位置的值，默认为0.0
        _minSlide.maximumValue = 100.0; // 滑块可以滑动到最小位置的值，默认为1.0
        _minSlide.tag = 100;
        [_minSlide addTarget:self action:@selector(clickSlide:) forControlEvents:UIControlEventValueChanged];
        [self.contentView addSubview:_minSlide];
        _minSlide.sd_layout.leftSpaceToView(_minLabel, 10).topSpaceToView(self.contentView, 10).rightSpaceToView(self.contentView, 15).heightIs(55);
        
        _maxslide = [[UISlider alloc]init];
        _maxslide.minimumValue = 16.0;  // 滑块可以滑动到最小位置的值，默认为0.0
        _maxslide.maximumValue = 100.0; // 滑块可以滑动到最小位置的值，默认为1.0
        _maxslide.tag = 101;
        [_maxslide addTarget:self action:@selector(clickSlide:) forControlEvents:UIControlEventValueChanged];
        [self.contentView addSubview:_maxslide];
        _maxslide.sd_layout.leftSpaceToView(_minLabel, 10).topSpaceToView(_minSlide, 10).rightSpaceToView(self.contentView, 15).heightIs(55);
    }
    return self;
}

-(void)requestMinVale:(NSString *)minValue maxValue:(NSString *)maxValue {
    _minSlide.value = [minValue integerValue];
    _maxslide.value = [maxValue integerValue];
}

-(void)clickSlide:(UISlider*)sender
{
    NSString *string = [NSString stringWithFormat:@"%.f",sender.value];
    
    if([_delegate respondsToSelector:@selector(clickSlideValueTag:withValue:)]) {
        [_delegate clickSlideValueTag:sender.tag withValue:string];
    }
}
@end
