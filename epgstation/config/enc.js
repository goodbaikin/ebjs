const spawn = require('child_process').spawn;
const execFile = require('child_process').execFile;
const ffprobe = process.env.FFPROBE;
const path = require('path');

const input = process.env.INPUT;
const output = process.env.OUTPUT;
const isDualMono = parseInt(process.env.AUDIOCOMPONENTTYPE, 10) == 2;
const input_name = path.basename(input, path.extname(output));
const output_name = path.basename(output);

const args = [
    '-server', process.env.EBJS_HOST,
    '-input', input_name,
    '-output', output_name,
    '-channelID', process.env.CHANNELID
];
if (isDualMono) {
    args.push('-isDualMono', 'true')
}

(async () => {
    // 進捗計算のために動画の長さを取得
    let duration = 0;

    //必要な変数
    let total_num = 0;
    let now_num = 0;
    let avisynth_flag = false;
    let percent = 0;
    let update_log_flag = false;
    let log = '';

    var env = Object.create(process.env);
    const child = spawn('ebjclient', args, { env: env });

    child.stdout.on('data', data => {
        let strbyline = String(data).split('\n');
        for (let i = 0; i < strbyline.length; i++) {
            let str = strbyline[i];
            switch (str) {
                case str.startsWith('chapter_exe') && str: {  //chapter_exe
                    const raw_chapter_exe_data = str.replace(/chapter_exe\s/, '');
                    switch (raw_chapter_exe_data) {
                        case raw_chapter_exe_data.startsWith('Creating') && str: {  //AviSynth+
                            const raw_avisynth_data = str;
                            const avisynth_reg = /Creating\slwi\sindex\sfile\s(\d+)%/;
                            total_num = 200;
                            now_num = Number(raw_avisynth_data.match(avisynth_reg)[1]);
                            now_num += avisynth_flag ? 100 : 0;
                            avisynth_flag = avisynth_flag ? true : now_num == 100 ? true : false;
                            update_log_flag = true;
                            log = `(1/4) AviSynth:Creating lwi index files`;
                            break;
                        }
                        case raw_chapter_exe_data.startsWith('\tVideo Frames') && raw_chapter_exe_data: {
                            //chapter_exeでの総フレーム数取得
                            const movie_frame_reg = /\tVideo\sFrames:\s(\d+)\s\[\d+\.\d+fps\]/;
                            total_num = Number(raw_chapter_exe_data.match(movie_frame_reg)[1]);
                            update_log_flag = true;
                            break;
                        }
                        case raw_chapter_exe_data.startsWith('mute') && raw_chapter_exe_data: {
                            //現在のフレーム数取得
                            const chapter_reg = /mute\s?\d+:\s(\d+)\s\-\s\d+フレーム/;
                            now_num = Number(raw_chapter_exe_data.match(chapter_reg)[1]);
                            update_log_flag = true;
                            break;
                        }
                        case raw_chapter_exe_data.startsWith('end') && raw_chapter_exe_data: {
                            //chapter_exeの終了検知
                            now_num = total_num;
                            update_log_flag = true;
                            break;
                        }
                        default: {
                            break;
                        }
                    }
                    log = `(2/4) Chapter_exe: ${now_num}/${total_num}`;
                    break;
                }

                case str.startsWith('logoframe') && str: { //logoframe
                    const raw_logoframe_data = str.replace(/logoframe\s/, '');
                    switch (raw_logoframe_data) {
                        case raw_logoframe_data.startsWith('checking') && raw_logoframe_data: {
                            const logoframe_reg = /checking\s*(\d+)\/(\d+)\sended./;
                            const logoframe = raw_logoframe_data.match(logoframe_reg);
                            if (logoframe == null) {
                                break;
                            }
                            now_num = Number(logoframe[1]);
                            total_num = Number(logoframe[2]);
                            update_log_flag = true;
                        }
                        default: {
                            break;
                        }
                    }
                    log = `(3/4) logoframe: ${now_num}/${total_num}`;
                    break;
                }

                case str.startsWith('ffprobe') && str: { //ffprobe
                    const raw_ffprobe_data = str.replace(/ffprobe\s/, '');
                    switch (raw_ffprobe_data) {
                        case raw_ffprobe_data.startsWith('duration') && raw_ffprobe_data: {
                            const duration_str = raw_ffprobe_data.split('=')[1].replace(/\r/g, '').trim();
                            duration = parseFloat(duration_str);
                        }
                        default: {
                            break;
                        }
                    }
                    break;
                }

                case str.startsWith('ffmpeg') && str: { //FFmpeg
                    // ffmpeg frame= 2847 fps=0.0 q=-1.0 Lsize=  216432kB time=00:01:35.64 bitrate=18537.1kbits/s speed= 222x
                    str = str.replace(/ffmpeg\s/, '');
                    if (!str.startsWith('frame')) {
                        continue;
                    }

                    const progress = {};
                    let tmp = (str + ' ').match(/[A-z]*=[A-z,0-9,\s,.,\/,:,-]* /g);
                    if (tmp === null) continue;
                    for (let j = 0; j < tmp.length; j++) {
                        progress[tmp[j].split('=')[0]] = tmp[j].split('=')[1].replace(/\r/g, '').trim();
                    }
                    progress['frame'] = parseInt(progress['frame']);
                    progress['fps'] = parseFloat(progress['fps']);
                    progress['q'] = parseFloat(progress['q']);

                    let current = 0;
                    if (progress.time == null) {
                        continue;
                    }
                    const times = progress.time.split(':');
                    for (let i = 0; i < times.length; i++) {
                        if (i == 0) {
                            current += parseFloat(times[i]) * 3600;
                        } else if (i == 1) {
                            current += parseFloat(times[i]) * 60;
                        } else if (i == 2) {
                            current += parseFloat(times[i]);
                        }
                    }

                    // 進捗率 1.0 で 100%
                    now_num = current;
                    total_num = duration;
                    log =
                        '(4/4) FFmpeg: ' +
                        //'frame= ' +
                        //progress.frame +
                        //' fps=' +
                        //progress.fps +
                        //' size=' +
                        //progress.size +
                        ' time=' +
                        progress.time +
                        //' bitrate=' +
                        //progress.bitrate +
                        ' speed=' +
                        progress.speed;
                    update_log_flag = true;
                    break;
                }

                default: { //進捗表示に必要ない出力データを流す
                    console.log(strbyline[i]);
                    break;
                }
            }
            percent = now_num / total_num;
            if (update_log_flag) console.log(JSON.stringify({ type: 'progress', percent: percent, log: log }));
            update_log_flag = false;
        }
    });

    child.on('error', err => {
        console.error(err);
        throw new Error(err);
    });

    process.on('SIGINT', () => {
        child.kill('SIGINT');
    });

    child.on('close', code => {
        //終了後にしたい処理があれば書く
    });
})();