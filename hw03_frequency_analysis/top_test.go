package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var textEn = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 Maecenas hendrerit est quis magna bibendum, sed malesuada erat dapibus. 
 Aliquam tempus metus elit, in pretium nisi rhoncus eu. 
 Vivamus consequat laoreet erat, et lacinia nulla faucibus quis. 
 Nullam eu tellus at risus ultrices pretium sit amet sed mauris. 
 Proin ultricies arcu diam, condimentum convallis elit congue a. 
 Fusce eleifend tempor sapien sed pharetra. 
 Maecenas blandit volutpat nunc in egestas. 

 Praesent laoreet consectetur tellus, sed feugiat neque scelerisque non. 
 Donec laoreet diam vel diam egestas consectetur. 
 Aliquam vel quam non nisl consectetur laoreet. 
 Donec consectetur nisi felis, ac cursus mauris efficitur quis.
 Praesent non ornare metus. Aliquam et suscipit mauris. Mauris varius non odio sed eleifend. 
 Phasellus sodales convallis mollis. Praesent diam mauris, euismod eu consequat vitae, laoreet quis elit. 
 In hac habitasse platea dictumst. Maecenas scelerisque, tortor a vestibulum pretium, 
 velit purus pharetra ex, vitae condimentum metus erat sit amet tellus. 
 Quisque ut eros auctor, eleifend elit eget, pretium metus. 
 Maecenas risus mauris, suscipit sed eleifend sed, vehicula non turpis.

 Vestibulum quis velit vitae erat vehicula dapibus non in tellus. 
 Fusce ornare imperdiet urna, ac lacinia dolor auctor sit amet. 
 Vivamus interdum lacinia augue id interdum. Quisque ut rutrum urna, eget ultricies mauris. 
 Fusce suscipit pharetra pharetra. Suspendisse lacinia efficitur turpis a hendrerit. 
 Ut odio odio, fermentum quis est ac, varius molestie mauris. Vestibulum id blandit tortor. 
 Donec tristique dui eu orci semper tempor. Nulla congue neque ut tempor egestas. 
 Duis imperdiet dui quis augue vehicula, non consectetur nisl imperdiet.

 Donec hendrerit nibh quis mauris cursus pellentesque. Aenean a scelerisque leo, sed semper tortor. 
 Sed eros diam, sodales sed bibendum at, vestibulum in ante. In quis rhoncus dui. 
 Donec auctor tincidunt tortor ut hendrerit. Sed bibendum fringilla orci, ut porttitor dolor dignissim hendrerit. 
 Quisque dapibus urna at sollicitudin imperdiet. Vivamus a fermentum purus.

 Curabitur a risus venenatis elit consequat feugiat. Maecenas vitae consectetur ante. 
 Sed congue blandit elit vitae tempus. Donec in tortor porttitor, suscipit turpis ac, posuere ligula. 
 Mauris faucibus nec justo posuere malesuada. Nulla facilities. 
 Cras elementum consectetur lacus, egestas egestas velit molestie nec. 
 Nulla vel gravida enim, id mollis lectus. Duis iaculis quam nunc, vitae facilisis libero faucibus nec.`

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test lorem", func(t *testing.T) {
		expected := []string{
			"sed",         // 12
			"mauris",      // 10
			"quis",        // 9
			"consectetur", // 8
			"elit",        // 7
			"in",          // 7
			"non",         // 7
			"a",           // 6
			"donec",       // 6
			"ut",          // 6
		}
		require.Equal(t, expected, Top10(textEn))
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}
